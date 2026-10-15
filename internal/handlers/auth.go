package handlers

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"time"
)

func (s *Server) validForm(w http.ResponseWriter, r *http.Request) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	if r.ParseForm() != nil || !equal(r.PostForm.Get("csrf"), s.csrf) {
		http.Error(w, "Invalid form token", http.StatusForbidden)
		return false
	}
	return true
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := struct{ CSRF, Error string }{CSRF: s.csrf}
	switch r.Method {
	case http.MethodGet:
		s.render(w, "login", data)
	case http.MethodPost:
		if !s.validForm(w, r) {
			return
		}
		// Bound failed attempts across all clients, including behind a reverse proxy.
		s.loginMu.Lock()
		now := time.Now()
		if !now.Before(s.loginRetryAt) {
			s.loginFailures = 0
			s.loginRetryAt = now.Add(time.Minute)
		}
		if s.loginFailures >= 5 {
			s.loginMu.Unlock()
			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusTooManyRequests)
			data.Error = "Too many attempts. Try again in one minute."
			s.render(w, "login", data)
			return
		}
		hash := sha256.Sum256([]byte(r.PostForm.Get("password")))
		ok := s.passwordSet && subtle.ConstantTimeCompare(hash[:], s.passwordHash[:]) == 1
		if !ok {
			s.loginFailures++
		}
		s.loginMu.Unlock()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			data.Error = "Incorrect password. Try again."
			s.render(w, "login", data)
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "divoom_session", Value: s.token, Path: "/", HttpOnly: true, SameSite: http.SameSiteStrictMode, Secure: r.TLS != nil})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
