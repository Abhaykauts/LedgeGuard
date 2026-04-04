package application

import (
	"errors"
	"time"

	"github.com/Abhaykauts/LedgeGuard/backend/internal/domain"
	"github.com/Abhaykauts/LedgeGuard/backend/pkg/auth"
)

type authService struct {
	userRepo      domain.UserRepository
	jwtSecret     string
	tokenDuration time.Duration
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(repo domain.UserRepository, secret string, duration time.Duration) AuthServiceInterface {
	return &authService{
		userRepo:      repo,
		jwtSecret:     secret,
		tokenDuration: duration,
	}
}

func (s *authService) Login(username, password string) (*AuthResponse, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !auth.CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := auth.GenerateToken(user.ID, string(user.Role), s.jwtSecret, s.tokenDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, err := auth.GenerateToken(user.ID, string(user.Role), s.jwtSecret, s.tokenDuration*24)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

func (s *authService) RefreshToken(token string) (*AuthResponse, error) {
	claims, err := auth.ValidateToken(token, s.jwtSecret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	accessToken, err := auth.GenerateToken(user.ID, string(user.Role), s.jwtSecret, s.tokenDuration)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: token, // Re-use refresh token or rotate
		User:         *user,
	}, nil
}
