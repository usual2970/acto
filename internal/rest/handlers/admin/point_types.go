package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/usual2970/acto/internal/rest/handlers"
	actoHttp "github.com/usual2970/acto/pkg/http"
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
	var req uc.PointTypeCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteError(w, 1000, "bad request")
		return
	}
	id, err := h.svc.Create(r.Context(), req)
	if err != nil {
		handlers.WriteDomainError(w, err)
		return
	}
	handlers.WriteSuccess(w, map[string]string{"id": id})
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
		handlers.WriteDomainError(w, err)
		return
	}
	handlers.WriteSuccess(w, res)
}

func (h *PointTypesHandler) Update(w http.ResponseWriter, r *http.Request) {
	// 从URL路径中获取积分类型名称
	vars := actoHttp.GetPathVars(r)
	name := vars["name"]
	if name == "" {
		handlers.WriteError(w, 1000, "missing name")
		return
	}

	// 解析部分更新字段
	var updates uc.PointTypeUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		handlers.WriteError(w, 1000, "bad request")
		return
	}

	if err := h.svc.Update(r.Context(), name, updates); err != nil {
		handlers.WriteDomainError(w, err)
		return
	}
	handlers.WriteSuccess(w, nil)
}

func (h *PointTypesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// 从URL路径中获取积分类型名称
	vars := actoHttp.GetPathVars(r)
	name := vars["name"]
	if name == "" {
		handlers.WriteError(w, 1000, "missing name")
		return
	}

	if err := h.svc.Delete(r.Context(), name); err != nil {
		handlers.WriteDomainError(w, err)
		return
	}
	handlers.WriteSuccess(w, nil)
}
