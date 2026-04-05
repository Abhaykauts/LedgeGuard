package http

import (
	"net/http"

	"github.com/Abhaykauts/LedgeGuard/backend/internal/application"
	"github.com/Abhaykauts/LedgeGuard/backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService application.AuthServiceInterface
}

func NewAuthHandler(service application.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{authService: service}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Login request"
// @Success 200 {object} application.AuthResponse
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.SendBadRequest(c, "invalid login request", err.Error())
		return
	}

	resp, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		errors.SendUnauthorized(c, "invalid credentials")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Refresh godoc
// @Summary Refresh access token
// @Description Use refresh token to get a new access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body refreshRequest true "Refresh request"
// @Success 200 {object} application.AuthResponse
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.SendBadRequest(c, "invalid refresh request", err.Error())
		return
	}

	resp, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		errors.SendUnauthorized(c, "invalid or expired refresh token")
		return
	}

	c.JSON(http.StatusOK, resp)
}
