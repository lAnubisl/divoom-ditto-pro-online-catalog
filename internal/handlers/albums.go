package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (s *Server) albumsPage(w http.ResponseWriter, r *http.Request) {
	data := pageData{CSRF: s.csrf, Automatic: s.automationState()}
	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()
	albums, err := s.galleryAPI.ListAlbums(ctx)
	if err != nil {
		log.Printf("albums: %v", err)
		data.Error = "Could not load albums. Check account access and reload the page."
	} else {
		data.Albums = albums
		if len(albums) > 0 {
			data.Album = albums[0].ID
			for _, a := range albums {
				if a.ID == data.Automatic.Album {
					data.Album = a.ID
				}
			}
		}
	}
	s.render(w, "albums", data)
}
func (s *Server) albumAnimations(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("album"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid album", 400)
		return
	}
	page := 1
	if value := r.URL.Query().Get("page"); value != "" {
		page, err = strconv.Atoi(value)
	}
	if err != nil || id <= 0 || page < 1 || page > 10000 {
		http.Error(w, "Invalid album/page", 400)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()
	data := pageData{Album: id, Page: page, Next: page + 1, Previous: page - 1, CSRF: s.csrf}
	result, err := s.galleryAPI.FetchAlbumPage(ctx, id, page)
	data.More = result.HasMore()
	if err == nil {
		data.Entries = result.Items
	}
	if err != nil {
		log.Printf("album %d: %v", id, err)
		data.Error = "Could not load the album. Try refreshing the list."
	}
	s.render(w, "animations", data)
}
