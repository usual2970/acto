package handlers

import (
	uc "acto/points"
	"encoding/json"
	"net/http"
	"strconv"
)

type RankingsHandler struct {
	ranking uc.RankingRepository
}

func NewRankingsHandler(r uc.RankingRepository) *RankingsHandler { return &RankingsHandler{ranking: r} }

func (h *RankingsHandler) Get(w http.ResponseWriter, r *http.Request) {
	// For simplicity, read query params
	pt := r.URL.Query().Get("pointTypeId")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 {
		limit = 100
	}
	start := int64(offset)
	stop := int64(offset + limit - 1)
	users, err := h.ranking.GetTop(r.Context(), pt, start, stop)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"items": users, "limit": limit, "offset": offset})
}
