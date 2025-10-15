package handlers

import (
	uc "acto/points"
	"encoding/json"
	"net/http"
	"strconv"
)

type RankingsHandler struct {
	ranking  uc.RankingRepository
	pointTyp uc.PointTypeRepository
}

func NewRankingsHandler(r uc.RankingRepository, pt uc.PointTypeRepository) *RankingsHandler {
	return &RankingsHandler{ranking: r, pointTyp: pt}
}

func (h *RankingsHandler) Get(w http.ResponseWriter, r *http.Request) {
	// Read pointTypeName and map to internal ID
	ptName := r.URL.Query().Get("pointTypeName")
	var ptID string
	if ptName != "" && h.pointTyp != nil {
		if pt, err := h.pointTyp.GetPointTypeByName(r.Context(), ptName); err == nil {
			ptID = pt.ID
		}
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 {
		limit = 100
	}
	start := int64(offset)
	stop := int64(offset + limit - 1)
	users, err := h.ranking.GetTop(r.Context(), ptID, start, stop)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"items": users, "limit": limit, "offset": offset})
}
