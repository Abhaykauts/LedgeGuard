package errors

import "github.com/gin-gonic/gin"

type AppError struct {
	Message string   `json:"error"`
	Code    int      `json:"-"`
	Details []string `json:"details,omitempty"`
}

func NewAppError(code int, message string, details ...string) *AppError {
	return &AppError{
		Message: message,
		Code:    code,
		Details: details,
	}
}

func (e *AppError) Send(c *gin.Context) {
	c.JSON(e.Code, e)
}

func SendBadRequest(c *gin.Context, message string, details ...string) {
	NewAppError(400, message, details...).Send(c)
}

func SendUnauthorized(c *gin.Context, message string) {
	NewAppError(401, message).Send(c)
}

func SendForbidden(c *gin.Context, message string) {
	NewAppError(403, message).Send(c)
}

func SendNotFound(c *gin.Context, message string) {
	NewAppError(404, message).Send(c)
}

func SendInternalError(c *gin.Context, message string) {
	NewAppError(500, message).Send(c)
}
