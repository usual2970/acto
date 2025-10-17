package handlers

import (
	"encoding/json"
	"net/http"

	uc "github.com/usual2970/acto/points"
)

type RedemptionsHandler struct{ svc *uc.RedemptionService }

func NewRedemptionsHandler(svc *uc.RedemptionService) *RedemptionsHandler {
	return &RedemptionsHandler{svc: svc}
}

func (h *RedemptionsHandler) Redeem(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   string `json:"userId"`
		RewardID string `json:"rewardId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.UserID == "" || req.RewardID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.svc.Redeem(r.Context(), req.UserID, req.RewardID); err != nil {
		writeDomainError(w, err, nil)
		return
	}
	w.WriteHeader(http.StatusOK)
}
