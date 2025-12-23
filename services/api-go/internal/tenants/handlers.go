package tenants

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler manages tenant-related HTTP endpoints
type Handler struct {
	service *Service
}

// NewHandler creates a new tenant handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateTenant handles POST /admin/tenants
func (h *Handler) CreateTenant(c *gin.Context) {
	var input CreateTenantInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	tenant, err := h.service.CreateTenant(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "create_tenant_failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"tenant": tenant,
		"message": "Tenant created successfully. Store the API key securely.",
	})
}

// RegenerateAPIKey handles POST /admin/tenants/:id/regenerate-key
func (h *Handler) RegenerateAPIKey(c *gin.Context) {
	tenantIDStr := c.Param("id")
	
	// Parse UUID
	tenantID, err := parseUUID(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_tenant_id",
			"message": "Tenant ID must be a valid UUID",
		})
		return
	}

	tenant, err := h.service.RegenerateAPIKey(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "regenerate_key_failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tenant":  tenant,
		"message": "API key regenerated successfully. Update your integrations.",
	})
}
