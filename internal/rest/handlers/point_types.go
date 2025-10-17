package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	uc "github.com/usual2970/acto/points"

	"github.com/gorilla/mux"
)

type PointTypesHandler struct {
	svc *uc.PointTypeService
}

func NewPointTypesHandler(svc *uc.PointTypeService) *PointTypesHandler {
	return &PointTypesHandler{svc: svc}
}

func (h *PointTypesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
		Description string `json:"description"`
	}
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
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	res, err := h.svc.List(r.Context(), limit, offset)
	if err != nil {
		writeDomainError(w, err, nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (h *PointTypesHandler) Update(w http.ResponseWriter, r *http.Request) {
	// 从URL路径中获取积分类型名称
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 解析部分更新字段
	var updates uc.UpdatePointTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.svc.Update(r.Context(), name, updates); err != nil {
		writeDomainError(w, err, nil)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *PointTypesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// 从URL路径中获取积分类型名称
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.svc.Delete(r.Context(), name); err != nil {
		writeDomainError(w, err, nil)
		return
	}
	w.WriteHeader(http.StatusOK)
}
