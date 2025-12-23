package persona

import (
	"net/http"

	"github.com/dennislee928/mighty-eagle/api-go/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler manages persona-related HTTP endpoints
type Handler struct {
	service *Service
}

// NewHandler creates a new persona handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateVerification handles POST /v1/persona/verifications
func (h *Handler) CreateVerification(c *gin.Context) {
	var input VerificationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	tenantID, _ := middleware.GetTenantID(c)

	verification, err := h.service.CreateVerification(c.Request.Context(), tenantID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "verification_error",
			"message": err.Error(),
		})
		return
	}

	statusCode := http.StatusCreated
	if verification.Status == "failed" {
		statusCode = http.StatusBadRequest // Or 422 Unprocessable Entity
	}

	c.JSON(statusCode, verification)
}

// GetVerification handles GET /v1/persona/verifications/:id
func (h *Handler) GetVerification(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "Verification ID must be a valid UUID",
		})
		return
	}

	tenantID, _ := middleware.GetTenantID(c)

	verification, err := h.service.GetVerification(c.Request.Context(), tenantID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "not_found",
			"message": "Verification not found",
		})
		return
	}

	c.JSON(http.StatusOK, verification)
}
