package admin

import (
	"encoding/json"
	"net/http"

	"github.com/usual2970/acto/internal/rest/handlers"
	uc "github.com/usual2970/acto/points"
)

type DistributionsHandler struct{ svc *uc.DistributionService }

func NewDistributionsHandler(svc *uc.DistributionService) *DistributionsHandler {
	return &DistributionsHandler{svc: svc}
}

func (h *DistributionsHandler) Execute(w http.ResponseWriter, r *http.Request) {
	var req uc.DistirbutionsExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteError(w, 1000, "bad request")
		return
	}
	if req.TopN <= 0 {
		req.TopN = 100
	}
	if err := h.svc.Execute(r.Context(), req); err != nil {
		handlers.WriteDomainError(w, err)
		return
	}
	handlers.WriteSuccess(w, nil)
}
