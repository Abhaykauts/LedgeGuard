package http

import (
	"net/http"
	"strconv"

	"github.com/Abhaykauts/LedgeGuard/backend/internal/application"
	"github.com/Abhaykauts/LedgeGuard/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type RecordHandler struct {
	service application.RecordServiceInterface
}

func NewRecordHandler(service application.RecordServiceInterface) *RecordHandler {
	return &RecordHandler{service: service}
}

// CreateRecord godoc
// @Summary Create a financial record
// @Description Add a new income or expense entry
// @Tags records
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param record body domain.Record true "Record object"
// @Success 201 {object} domain.Record
// @Router /api/records [post]
func (h *RecordHandler) CreateRecord(c *gin.Context) {
	var record domain.Record
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	record.CreatedBy = userID.(uint)

	if err := h.service.CreateRecord(&record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, record)
}

// ListRecords godoc
// @Summary List financial records
// @Description Get all records with optional filters
// @Tags records
// @Security BearerAuth
// @Produce json
// @Param type query string false "Filter by type (INCOME/EXPENSE)"
// @Param category query string false "Filter by category"
// @Success 200 {array} domain.Record
// @Router /api/records [get]
func (h *RecordHandler) ListRecords(c *gin.Context) {
	var filter domain.RecordFilter
	
	if t := c.Query("type"); t != "" {
		recType := domain.RecordType(t)
		filter.Type = &recType
	}
	if cat := c.Query("category"); cat != "" {
		filter.Category = &cat
	}

	records, err := h.service.ListRecords(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

// DeleteRecord godoc
func (h *RecordHandler) DeleteRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := h.service.DeleteRecord(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
