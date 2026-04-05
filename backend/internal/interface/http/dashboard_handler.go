package http

import (
	"net/http"

	"github.com/Abhaykauts/LedgeGuard/backend/internal/application"
	"github.com/Abhaykauts/LedgeGuard/backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	service application.DashboardServiceInterface
}

func NewDashboardHandler(service application.DashboardServiceInterface) *DashboardHandler {
	return &DashboardHandler{service: service}
}

// GetSummary godoc
// @Summary Get dashboard summary
// @Description Get total income, expenses, and category-wise totals
// @Tags dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} application.DashboardSummary
// @Failure 500 {object} errors.AppError
// @Router /dashboard/summary [get]
func (h *DashboardHandler) GetSummary(c *gin.Context) {
	summary, err := h.service.GetSummary()
	if err != nil {
		errors.SendInternalError(c, "failed to get dashboard summary")
		return
	}

	c.JSON(http.StatusOK, summary)
}
