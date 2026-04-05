package http

import (
	"net/http"
	"strconv"

	"github.com/Abhaykauts/LedgeGuard/backend/internal/domain"
	"github.com/Abhaykauts/LedgeGuard/backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repo domain.UserRepository
}

func NewUserHandler(repo domain.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// @Summary List all users
// @Description Get a list of all users. Admin only.
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.User
// @Failure 500 {object} errors.AppError
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.repo.List()
	if err != nil {
		errors.SendInternalError(c, "failed to fetch users")
		return
	}
	c.JSON(http.StatusOK, users)
}

// @Summary Create a new user
// @Description Create a new system user. Admin only.
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user body domain.User true "User object"
// @Success 201 {object} domain.User
// @Failure 400 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		errors.SendBadRequest(c, "invalid user data", err.Error())
		return
	}

	if err := h.repo.Create(&user); err != nil {
		errors.SendInternalError(c, "failed to create user")
		return
	}

	c.JSON(http.StatusCreated, user)
}

// @Summary Update a user
// @Description Update an existing user. Admin only.
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body domain.User true "User object"
// @Success 200 {object} domain.User
// @Failure 400 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		errors.SendBadRequest(c, "invalid user data", err.Error())
		return
	}
	user.ID = uint(id)

	if err := h.repo.Update(&user); err != nil {
		errors.SendInternalError(c, "failed to update user")
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Delete a user (Soft Delete)
// @Description Deactivate a user. Admin only.
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 500 {object} errors.AppError
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.repo.Delete(uint(id)); err != nil {
		errors.SendInternalError(c, "failed to delete user")
		return
	}
	c.Status(http.StatusNoContent)
}
