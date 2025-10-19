package api

import (
	"encoding/json"
	"net/http"

	"github.com/usual2970/acto/internal/rest/handlers"
	uc "github.com/usual2970/acto/points"
)

type RedemptionsHandler struct{ svc *uc.RedemptionService }

func NewRedemptionsHandler(svc *uc.RedemptionService) *RedemptionsHandler {
	return &RedemptionsHandler{svc: svc}
}

func (h *RedemptionsHandler) Redeem(w http.ResponseWriter, r *http.Request) {
	var req uc.RedemptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteError(w, 1000, "bad request")
		return
	}
	if req.UserID == "" || req.RewardID == "" {
		handlers.WriteError(w, 1000, "missing userId or rewardId")
		return
	}
	if err := h.svc.Redeem(r.Context(), req); err != nil {
		handlers.WriteDomainError(w, err)
		return
	}
	handlers.WriteSuccess(w, nil)
}
