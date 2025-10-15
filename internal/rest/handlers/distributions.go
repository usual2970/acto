package handlers

import (
	uc "acto/points"
	"encoding/json"
	"net/http"
)

type DistributionsHandler struct{ svc *uc.DistributionService }

func NewDistributionsHandler(svc *uc.DistributionService) *DistributionsHandler {
	return &DistributionsHandler{svc: svc}
}

func (h *DistributionsHandler) Execute(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PointTypeName string
		TopN          int
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.TopN <= 0 {
		req.TopN = 100
	}
	if err := h.svc.Execute(r.Context(), req.PointTypeName, req.TopN); err != nil {
		writeDomainError(w, err, nil)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
