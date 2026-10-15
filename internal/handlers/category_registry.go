package handlers

import "divoom-gallery-client/internal/gallery"

// The category list is immutable for the lifetime of the server.
func (s *Server) categoryList() []gallery.Category {
	return append([]gallery.Category(nil), s.categories...)
}
func (s *Server) hasCategory(id int) bool {
	for _, c := range s.categories {
		if c.ID == id {
			return true
		}
	}
	return false
}
