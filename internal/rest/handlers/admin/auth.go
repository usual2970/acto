package admin

import (
	"encoding/json"
	"net/http"

	uc "github.com/usual2970/acto/auth"

	"github.com/usual2970/acto/internal/rest/handlers"
)

type AuthHandler struct {
	svc *uc.AuthService
}

func NewAuthHandler(svc *uc.AuthService) *AuthHandler {
	return &AuthHandler{
		svc: svc,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var req uc.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteError(w, 1000, "bad request")
		return
	}
	token, err := h.svc.Authenticate(req)
	if err != nil {
		handlers.WriteError(w, 1000, err.Error())
		return
	}
	handlers.WriteSuccess(w, map[string]string{"token": token})
}
