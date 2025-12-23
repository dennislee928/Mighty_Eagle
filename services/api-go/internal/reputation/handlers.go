package reputation

import (
	"net/http"

	"github.com/dennislee928/mighty-eagle/api-go/internal/middleware"
	"github.com/gin-gonic/gin"
)

// Handler manages reputation-related HTTP endpoints
type Handler struct {
	service *Service
}

// NewHandler creates a new reputation handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetReputation handles GET /v1/reputation/:subject
func (h *Handler) GetReputation(c *gin.Context) {
	subjectID := c.Param("subject")
	if subjectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Subject ID is required",
		})
		return
	}

	tenantID, _ := middleware.GetTenantID(c)

	result, err := h.service.GetReputation(c.Request.Context(), tenantID, subjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "reputation_error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
