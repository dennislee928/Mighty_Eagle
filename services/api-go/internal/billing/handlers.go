package billing

import (
	"net/http"

	"github.com/dennislee928/mighty-eagle/api-go/internal/middleware"
	"github.com/gin-gonic/gin"
)

// Handler manages billing-related HTTP endpoints
type Handler struct {
	service *Service
}

// NewHandler creates a new billing handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetUsage handles GET /v1/billing/usage
func (h *Handler) GetUsage(c *gin.Context) {
	tenant, ok := middleware.GetTenant(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	usage, err := h.service.GetCurrentUsage(c.Request.Context(), *tenant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "usage_check_failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"period": "current_month",
		"tier":   tenant.Tier,
		"metrics": usage,
	})
}
