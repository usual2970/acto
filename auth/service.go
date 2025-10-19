package auth

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/usual2970/acto/internal/config"
)

// AuthService authenticates users against credentials configured via env vars
// and issues a signed JWT on success.
//
// Required env vars:
// - AUTH_USERNAME: expected username
// - AUTH_PASSWORD: expected password
// - JWT_SECRET: HMAC secret used to sign tokens (HS256)
// Optional env vars:
// - JWT_TTL: token lifetime as Go duration (e.g. "1h", "30m"); defaults to 1h
// - JWT_ISSUER: token issuer; defaults to "acto-auth"
type AuthService struct {
	secret       []byte
	issuer       string
	ttl          time.Duration
	expectedUser string
	expectedPass string
}

// NewAuthService creates an AuthService from process config
func NewAuthService() *AuthService {
	cfg := config.Load()
	return NewAuthServiceWithConfig(cfg)
}

// NewAuthServiceWithConfig creates an AuthService from provided config
func NewAuthServiceWithConfig(cfg config.Config) *AuthService {
	ttl := time.Hour
	if cfg.JWTTTL != "" {
		if d, err := time.ParseDuration(cfg.JWTTTL); err == nil {
			ttl = d
		}
	}
	return &AuthService{
		secret:       []byte(cfg.JWTSecret),
		issuer:       cfg.JWTIssuer,
		ttl:          ttl,
		expectedUser: cfg.AuthUsername,
		expectedPass: cfg.AuthPassword,
	}
}

// Authenticate validates provided credentials and returns a signed JWT token string.
func (s *AuthService) Authenticate(req AuthRequest) (string, error) {
	if s.expectedUser == "" || s.expectedPass == "" {
		return "", errors.New("auth credentials not configured")
	}
	if req.Username != s.expectedUser || req.Password != s.expectedPass {
		return "", errors.New("invalid credentials")
	}
	if len(s.secret) == 0 {
		return "", errors.New("jwt secret not configured")
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  req.Username,
		"role": "admin", // by design, env-backed login is for admin console
		"iss":  s.issuer,
		"iat":  now.Unix(),
		"exp":  now.Add(s.ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", err
	}
	return signed, nil
}
