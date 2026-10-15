package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (s *Server) categoriesPage(w http.ResponseWriter, r *http.Request) {
	categories := s.categoryList()
	category := s.category
	if !s.hasCategory(category) {
		category = categories[0].ID
	}
	s.render(w, "categories", pageData{Categories: categories, Category: category, CSRF: s.csrf, Automatic: s.automationState()})
}
func (s *Server) categoryAnimations(w http.ResponseWriter, r *http.Request) {
	category, err := strconv.Atoi(r.URL.Query().Get("category"))
	if err != nil {
		http.Error(w, "Invalid category", 400)
		return
	}
	valid := s.hasCategory(category)
	page := 1
	if value := r.URL.Query().Get("page"); value != "" {
		page, err = strconv.Atoi(value)
	}
	if err != nil || !valid || page < 1 || page > 10000 {
		http.Error(w, "Invalid category or page", 400)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()
	result, err := s.galleryAPI.FetchCategoryPage(ctx, category, page)
	data := pageData{Category: category, Page: page, Next: page + 1, Previous: page - 1, More: result.HasMore(), CSRF: s.csrf}
	if err == nil {
		data.Entries = result.Items
	}
	if err != nil {
		log.Printf("category animations: %v", err)
		data.Error = "Could not load the category. Try refreshing the list."
	}
	s.render(w, "animations", data)
}
