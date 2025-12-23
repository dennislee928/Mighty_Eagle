package audit

import (
	"net/http"
	"time"

	"github.com/dennislee928/mighty-eagle/api-go/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler manages audit-related HTTP endpoints
type Handler struct {
	exporter *Exporter
	// logger *Logger // If we need direct log access
}

// NewHandler creates a new audit handler
func NewHandler(exporter *Exporter) *Handler {
	return &Handler{exporter: exporter}
}

type CreateExportInput struct {
	Format    string    `json:"format" binding:"required"` // csv or json
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

// CreateExport handles POST /v1/audit/exports
func (h *Handler) CreateExport(c *gin.Context) {
	var input CreateExportInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	tenantID, _ := middleware.GetTenantID(c)

	job, err := h.exporter.CreateExportJob(c.Request.Context(), tenantID, input.Format, input.StartDate, input.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "export_failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, job)
}

// GetExport handles GET /v1/audit/exports/:id
func (h *Handler) GetExport(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "Job ID must be a valid UUID",
		})
		return
	}

	tenantID, _ := middleware.GetTenantID(c)

	job, err := h.exporter.GetExportJob(c.Request.Context(), tenantID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "not_found",
			"message": "Export job not found",
		})
		return
	}

	c.JSON(http.StatusOK, job)
}
