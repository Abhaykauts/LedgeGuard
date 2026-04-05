package http

import (
	"net/http"
	"strconv"

	"github.com/Abhaykauts/LedgeGuard/backend/internal/application"
	"github.com/Abhaykauts/LedgeGuard/backend/internal/domain"
	"github.com/Abhaykauts/LedgeGuard/backend/pkg/errors"
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
// @Description Add a new income or expense entry. Admin only.
// @Tags records
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param record body domain.Record true "Record object"
// @Success 201 {object} domain.Record
// @Failure 400 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /records [post]
func (h *RecordHandler) CreateRecord(c *gin.Context) {
	var record domain.Record
	if err := c.ShouldBindJSON(&record); err != nil {
		errors.SendBadRequest(c, "invalid record data", err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	record.CreatedBy = userID.(uint)

	if err := h.service.CreateRecord(&record); err != nil {
		errors.SendInternalError(c, "failed to create record")
		return
	}

	c.JSON(http.StatusCreated, record)
}

// ListRecords godoc
// @Summary List financial records
// @Description Get records with advanced filtering and pagination.
// @Tags records
// @Security BearerAuth
// @Produce json
// @Param start_date query string false "Filter start (YYYY-MM-DD)"
// @Param end_date query string false "Filter end (YYYY-MM-DD)"
// @Param type query string false "Filter by type (INCOME/EXPENSE)"
// @Param category query string false "Filter by category"
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Page size (default 10)"
// @Success 200 {array} domain.Record
// @Failure 400 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /records [get]
func (h *RecordHandler) ListRecords(c *gin.Context) {
	var filter domain.RecordFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		errors.SendBadRequest(c, "invalid filter parameters")
		return
	}

	records, err := h.service.ListRecords(filter)
	if err != nil {
		errors.SendInternalError(c, "failed to fetch records")
		return
	}

	c.JSON(http.StatusOK, records)
}

// UpdateRecord godoc
// @Summary Update a financial record
// @Description Modify an existing record. Admin only.
// @Tags records
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Record ID"
// @Param record body domain.Record true "Record object"
// @Success 200 {object} domain.Record
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /records/{id} [put]
func (h *RecordHandler) UpdateRecord(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// Check if record exists
	existing, err := h.service.GetRecord(uint(id))
	if err != nil || existing == nil {
		errors.SendNotFound(c, "record not found")
		return
	}

	var record domain.Record
	if err := c.ShouldBindJSON(&record); err != nil {
		errors.SendBadRequest(c, "invalid input data")
		return
	}
	record.ID = uint(id)

	if err := h.service.UpdateRecord(&record); err != nil {
		errors.SendInternalError(c, "failed to update record")
		return
	}

	c.JSON(http.StatusOK, record)
}

// DeleteRecord godoc
// @Summary Delete a financial record
// @Description Soft delete a record. Admin only.
// @Tags records
// @Security BearerAuth
// @Param id path int true "Record ID"
// @Success 204 "No Content"
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /records/{id} [delete]
func (h *RecordHandler) DeleteRecord(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// Check existence before delete
	existing, err := h.service.GetRecord(uint(id))
	if err != nil || existing == nil {
		errors.SendNotFound(c, "record not found")
		return
	}

	if err := h.service.DeleteRecord(uint(id)); err != nil {
		errors.SendInternalError(c, "failed to delete record")
		return
	}
	c.Status(http.StatusNoContent)
}
