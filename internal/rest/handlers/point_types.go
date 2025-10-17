package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/usual2970/acto/pkg"
	uc "github.com/usual2970/acto/points"
	// path vars are read via request context to stay framework-agnostic
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
		WriteError(w, 1000, "bad request")
		return
	}
	id, err := h.svc.Create(r.Context(), req.Name, req.DisplayName, req.Description)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	WriteSuccess(w, map[string]string{"id": id})
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
		writeDomainError(w, err)
		return
	}
	WriteSuccess(w, res)
}

func (h *PointTypesHandler) Update(w http.ResponseWriter, r *http.Request) {
	// 从URL路径中获取积分类型名称
	vars := pkg.GetPathVars(r)
	name := vars["name"]
	if name == "" {
		WriteError(w, 1000, "missing name")
		return
	}

	// 解析部分更新字段
	var updates uc.UpdatePointTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		WriteError(w, 1000, "bad request")
		return
	}

	if err := h.svc.Update(r.Context(), name, updates); err != nil {
		writeDomainError(w, err)
		return
	}
	WriteSuccess(w, nil)
}

func (h *PointTypesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// 从URL路径中获取积分类型名称
	vars := pkg.GetPathVars(r)
	name := vars["name"]
	if name == "" {
		WriteError(w, 1000, "missing name")
		return
	}

	if err := h.svc.Delete(r.Context(), name); err != nil {
		writeDomainError(w, err)
		return
	}
	WriteSuccess(w, nil)
}
