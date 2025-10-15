package handlers

import (
	d "acto/domain/points"
	uc "acto/points"
	"encoding/json"
	"net/http"
	"strconv"
)

type PointTypesHandler struct {
	svc *uc.PointTypeService
}

func NewPointTypesHandler(svc *uc.PointTypeService) *PointTypesHandler {
	return &PointTypesHandler{svc: svc}
}

func (h *PointTypesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct{ Name, DisplayName, Description string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := h.svc.Create(r.Context(), req.Name, req.DisplayName, req.Description)
	if err != nil {
		writeDomainError(w, err, nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (h *PointTypesHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	res, err := h.svc.List(r.Context(), limit, offset)
	if err != nil {
		writeDomainError(w, err, nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (h *PointTypesHandler) Update(w http.ResponseWriter, r *http.Request) {
	var pt d.PointType
	if err := json.NewDecoder(r.Body).Decode(&pt); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.svc.Update(r.Context(), pt); err != nil {
		writeDomainError(w, err, nil)
		return
	}
	w.WriteHeader(http.StatusOK)
}
