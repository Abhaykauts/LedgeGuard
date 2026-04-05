package http

import (
	"net/http"

	"github.com/Abhaykauts/LedgeGuard/backend/internal/application"
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
// @Router /api/dashboard/summary [get]
func (h *DashboardHandler) GetSummary(c *gin.Context) {
	summary, err := h.service.GetSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}
