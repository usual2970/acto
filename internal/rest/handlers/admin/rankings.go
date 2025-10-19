package admin

import (
	"net/http"
	"strconv"

	"github.com/usual2970/acto/internal/rest/handlers"
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
		handlers.WriteDomainError(w, err)
		return
	}
	handlers.WriteSuccess(w, map[string]any{"items": users, "limit": limit, "offset": offset})
}
