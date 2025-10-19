package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/usual2970/acto/internal/rest/handlers"
	actoHttp "github.com/usual2970/acto/pkg/http"
	uc "github.com/usual2970/acto/points"
	// path vars are read via request context to stay framework-agnostic
)

type BalancesHandler struct {
	svc *uc.BalanceService
}

func NewBalancesHandler(svc *uc.BalanceService) *BalancesHandler { return &BalancesHandler{svc: svc} }

func (h *BalancesHandler) Credit(w http.ResponseWriter, r *http.Request) {

	var req uc.BalanceCreditRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteError(w, 1000, "bad request")
		return
	}
	if err := h.svc.Credit(r.Context(), req); err != nil {
		handlers.WriteDomainError(w, err)
		return
	}
	handlers.WriteSuccess(w, nil)
}

func (h *BalancesHandler) Debit(w http.ResponseWriter, r *http.Request) {
	var req uc.BalanceDebitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteError(w, 1000, "bad request")
		return
	}
	if err := h.svc.Debit(r.Context(), req); err != nil {
		handlers.WriteDomainError(w, err)
		return
	}
	handlers.WriteSuccess(w, nil)
}

func (h *BalancesHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	// Parse params
	vars := actoHttp.GetPathVars(r)
	userID := vars["userId"]

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	pointTypeName := r.URL.Query().Get("pointTypeName")
	op := r.URL.Query().Get("op")
	startTime, _ := strconv.ParseInt(r.URL.Query().Get("startTime"), 10, 64)
	endTime, _ := strconv.ParseInt(r.URL.Query().Get("endTime"), 10, 64)
	items, total, err := h.svc.ListTransactions(r.Context(), userID, pointTypeName, op, startTime, endTime, limit, offset)
	if err != nil {
		handlers.WriteDomainError(w, err)
		return
	}
	handlers.WriteSuccess(w, map[string]any{"items": items, "total": total, "limit": limit, "offset": offset})
}
