package application

import "github.com/Abhaykauts/LedgeGuard/backend/internal/domain"

// AuthResponse contains the session tokens
type AuthResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         domain.User `json:"user"`
}

// AuthServiceInterface defines the authentication use cases
type AuthServiceInterface interface {
	Login(username, password string) (*AuthResponse, error)
	RefreshToken(token string) (*AuthResponse, error)
}
