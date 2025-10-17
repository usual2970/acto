package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	uc "github.com/usual2970/acto/points"

	"github.com/gorilla/mux"
)

type BalancesHandler struct {
	svc *uc.BalanceService
}

func NewBalancesHandler(svc *uc.BalanceService) *BalancesHandler { return &BalancesHandler{svc: svc} }

func (h *BalancesHandler) Credit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID        string `json:"userId"`
		PointTypeName string `json:"pointTypeName"`
		Reason        string `json:"reason"`
		Amount        int64  `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.svc.Credit(r.Context(), req.UserID, req.PointTypeName, req.Reason, req.Amount); err != nil {
		writeDomainError(w, err, nil)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *BalancesHandler) Debit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID        string `json:"userId"`
		PointTypeName string `json:"pointTypeName"`
		Reason        string `json:"reason"`
		Amount        int64  `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.svc.Debit(r.Context(), req.UserID, req.PointTypeName, req.Reason, req.Amount); err != nil {
		writeDomainError(w, err, nil)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *BalancesHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	// Parse params
	vars := mux.Vars(r)
	userID := vars["userId"]

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	pointTypeName := r.URL.Query().Get("pointTypeName")
	op := r.URL.Query().Get("op")
	startTime, _ := strconv.ParseInt(r.URL.Query().Get("startTime"), 10, 64)
	endTime, _ := strconv.ParseInt(r.URL.Query().Get("endTime"), 10, 64)
	items, total, err := h.svc.ListTransactions(r.Context(), userID, pointTypeName, op, startTime, endTime, limit, offset)
	if err != nil {
		writeDomainError(w, err, nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"items": items, "total": total, "limit": limit, "offset": offset})
}
