package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/frankfoo/doraemon/internal/classify"
	"github.com/frankfoo/doraemon/internal/store"
)

type entryRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Tags     string `json:"tags"`
}

type paginatedResponse struct {
	Entries    []store.Entry `json:"entries"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PerPage    int           `json:"per_page"`
	TotalPages int           `json:"total_pages"`
}

func (s *Server) handleListEntries(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	total, err := s.store.Count(category)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to count entries")
		return
	}

	offset := (page - 1) * perPage
	entries, err := s.store.ListPaginated(category, offset, perPage)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list entries")
		return
	}
	if entries == nil {
		entries = []store.Entry{}
	}

	totalPages := (total + perPage - 1) / perPage
	writeJSON(w, http.StatusOK, paginatedResponse{
		Entries:    entries,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	})
}

func (s *Server) handleGetEntry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	entry, err := s.store.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get entry")
		return
	}
	if entry == nil {
		writeError(w, http.StatusNotFound, "Entry not found")
		return
	}

	writeJSON(w, http.StatusOK, entry)
}

func (s *Server) handleCreateEntry(w http.ResponseWriter, r *http.Request) {
	var req entryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.Title == "" || req.Content == "" {
		writeError(w, http.StatusBadRequest, "Title and content are required")
		return
	}

	category := req.Category
	tags := req.Tags
	if category == "" || tags == "" {
		result := classify.Classify(req.Title, req.Content)
		if category == "" {
			category = result.Category
		}
		if tags == "" {
			tags = classify.FormatTags(result.Tags)
		}
	}

	entry := &store.Entry{
		Title:    req.Title,
		Content:  req.Content,
		Category: category,
		Tags:     tags,
	}

	id, err := s.store.Save(entry)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create entry")
		return
	}

	created, err := s.store.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get created entry")
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) handleUpdateEntry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	existing, err := s.store.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get entry")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "Entry not found")
		return
	}

	var req entryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title != "" {
		existing.Title = req.Title
	}
	if req.Content != "" {
		existing.Content = req.Content
	}
	if req.Category != "" {
		existing.Category = req.Category
	}
	if req.Tags != "" {
		existing.Tags = req.Tags
	}

	if err := s.store.Update(existing); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update entry")
		return
	}

	writeJSON(w, http.StatusOK, existing)
}

func (s *Server) handleDeleteEntry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	existing, err := s.store.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get entry")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "Entry not found")
		return
	}

	if err := s.store.Delete(id); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete entry")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Entry deleted"})
}

func (s *Server) handleSearchEntries(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		writeError(w, http.StatusBadRequest, "Query parameter 'q' is required")
		return
	}

	query := store.BuildFTSQuery(q)
	entries, err := s.store.Search(query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to search entries")
		return
	}
	if entries == nil {
		entries = []store.Entry{}
	}

	writeJSON(w, http.StatusOK, entries)
}
