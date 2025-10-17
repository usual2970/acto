package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	uc "github.com/usual2970/acto/points"
)

type RankingsHandler struct {
	svc uc.RankingsService
}

func NewRankingsHandler(svc uc.RankingsService) *RankingsHandler {
	return &RankingsHandler{svc: svc}
}

func (h *RankingsHandler) Get(w http.ResponseWriter, r *http.Request) {
	ptName := r.URL.Query().Get("pointTypeName")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	users, err := h.svc.GetTop(r.Context(), ptName, limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"items": users, "limit": limit, "offset": offset})
}
