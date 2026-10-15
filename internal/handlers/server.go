package handlers

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"divoom-gallery-client/internal/gallery"
	"divoom-gallery-client/internal/playback"
	"divoom-gallery-client/internal/rotation"
	"divoom-gallery-client/internal/settings"
	"divoom-gallery-client/internal/source"
	"divoom-gallery-client/web"
)

// Server handles HTTP routing, authentication and template rendering.
type Server struct {
	galleryAPI    *gallery.Client
	source        *source.Adapter
	playback      *playback.Controller
	templates     *template.Template
	token, csrf   string
	itemKey       string
	passwordHash  [32]byte
	passwordSet   bool
	loginMu       sync.Mutex
	loginFailures int
	loginRetryAt  time.Time
	category      int
	settings      *settings.Store
	runner        *rotation.Runner
	categories    []gallery.Category
}

func New(g *source.Adapter, galleryAPI *gallery.Client, commands *playback.Controller, prefs *settings.Store, runner *rotation.Runner, category int, password string) *Server {
	s := &Server{source: g, galleryAPI: galleryAPI, playback: commands, settings: prefs, runner: runner, category: category,
		passwordHash: sha256.Sum256([]byte(password)), passwordSet: password != "",
		token: secret(), csrf: secret(), itemKey: secret()}
	s.templates = template.Must(template.New("pages").Funcs(template.FuncMap{"itemToken": s.itemToken}).ParseFS(web.Assets, "templates/*.html"))
	s.categories = galleryAPI.ListCategories()
	if !s.hasCategory(category) {
		s.categories = append(s.categories, gallery.Category{ID: category, Name: strconv.Itoa(category)})
	}
	// Include a saved numeric source without changing the browsing default.
	if saved, err := prefs.Read(); err == nil && saved.Source == "category" && !s.hasCategory(int(saved.SourceID)) {
		s.categories = append(s.categories, gallery.Category{ID: int(saved.SourceID), Name: strconv.FormatInt(saved.SourceID, 10)})
	}
	return s
}
func secret() string {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b[:])
}
func equal(a, b string) bool { return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1 }

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	// Categories.
	mux.HandleFunc("GET /{$}", s.categoriesPage)
	mux.HandleFunc("GET /categories", s.categoriesPage)
	mux.HandleFunc("GET /categories/animations", s.categoryAnimations)

	// Albums.
	mux.HandleFunc("GET /albums", s.albumsPage)
	mux.HandleFunc("GET /albums/animations", s.albumAnimations)

	// Autoplay settings and status.
	mux.HandleFunc("GET /autoplay", s.settingsPage)
	mux.HandleFunc("GET /automation", s.automationStatus)
	mux.HandleFunc("GET /automation/widget", s.automationWidget)
	mux.HandleFunc("POST /automation/category", s.setAutoplayCategory)
	mux.HandleFunc("POST /automation/album", s.setAutoplayAlbum)
	mux.HandleFunc("POST /automation/resume", s.resumeAutoplay)
	mux.HandleFunc("POST /automation/pause", s.pauseAutoplay)
	mux.HandleFunc("POST /automation/timing", s.setAutoplayTiming)

	// Clock.
	mux.HandleFunc("GET /clock", s.clockPage)
	mux.HandleFunc("POST /clock/settings", s.applyClockSettings)
	mux.HandleFunc("POST /clock", s.restoreClock)

	// Animation preview and playback.
	mux.HandleFunc("GET /animations/{id}", s.preview)
	mux.HandleFunc("POST /send/{id}", s.send)

	// Static assets.
	mux.Handle("GET /static/", http.FileServerFS(web.Assets))

	return securityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" {
			s.login(w, r)
			return
		}
		if r.Method == "GET" && r.URL.Path == "/static/css/style.css" {
			mux.ServeHTTP(w, r)
			return
		}
		cookie, err := r.Cookie("divoom_session")
		if err != nil || !equal(cookie.Value, s.token) {
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/login")
			} else if r.Method == "GET" || r.Method == "HEAD" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			http.Error(w, "Sign in to continue.", http.StatusUnauthorized)
			return
		}
		if r.Method == "POST" {
			if !s.validForm(w, r) {
				return
			}
		}
		mux.ServeHTTP(w, r)
	}))
}

func (s *Server) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.templates.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("render %s: %v", name, err)
	}
}
func animationID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid animation ID")
	}
	return id, nil
}
func (s *Server) preview(w http.ResponseWriter, r *http.Request) {
	item, err := s.animationItem(r)
	if err != nil {
		http.Error(w, "Invalid animation item", 400)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()
	data, err := s.source.Animation(ctx, item)
	if err != nil {
		log.Printf("preview animation %d (type=%d size=%d): %v", item.ID, item.Type, item.Size, err)
		w.Header().Set("Content-Type", "image/svg+xml")
		fmt.Fprint(w, `<svg xmlns="http://www.w3.org/2000/svg" width="160" height="160"><rect width="160" height="160" fill="#202431"/><text x="80" y="78" text-anchor="middle" fill="#aab2c5" font-size="12">No preview</text><text x="80" y="100" text-anchor="middle" fill="#aab2c5" font-size="10">Unsupported format</text></svg>`)
		return
	}
	w.Header().Set("Content-Type", "image/gif")
	w.Write(data)
}
func (s *Server) status(w http.ResponseWriter, message string) {
	s.render(w, "status", message)
	s.render(w, "automation-oob", s.automationState())
}
func (s *Server) send(w http.ResponseWriter, r *http.Request) {
	item, err := s.animationItem(r)
	if err != nil {
		http.Error(w, "Invalid animation item", 400)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 480*time.Second)
	defer cancel()
	if err := s.playback.Send(ctx, item); err != nil {
		log.Printf("send animation %d: %v", item.ID, err)
		s.status(w, commandMessage(err, "Send failed. The device state is unknown; you can restore the clock manually."))
		return
	}
	s.status(w, fmt.Sprintf("Animation #%d sent: the server confirmed completion.", item.ID))
}
