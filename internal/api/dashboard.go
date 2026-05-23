package api

import (
	"net/http"
)

type dashboardResponse struct {
	TotalEntries  int            `json:"total_entries"`
	ByCategory    map[string]int `json:"by_category"`
	RecentEntries any            `json:"recent_entries"`
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	total, err := s.store.Count("")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get total count")
		return
	}

	byCategory, err := s.store.CountByCategory()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get category counts")
		return
	}

	recent, err := s.store.RecentEntries(5)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get recent entries")
		return
	}

	writeJSON(w, http.StatusOK, dashboardResponse{
		TotalEntries:  total,
		ByCategory:    byCategory,
		RecentEntries: recent,
	})
}
