package auth

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
