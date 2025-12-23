package consent

import (
	"net/http"

	"github.com/dennislee928/mighty-eagle/api-go/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler manages consent-related HTTP endpoints
type Handler struct {
	service *Service
}

// NewHandler creates a new consent handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateToken handles POST /v1/consent/tokens
func (h *Handler) CreateToken(c *gin.Context) {
	var input CreateTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	tenantID, _ := middleware.GetTenantID(c)
	// Get tenant secret for signing receipts
	// In M0 we added APISecretHash to tenant model, but we don't store the plain secret usually?
	// Oh, wait, we do generate it. We probably need a way to get the signing secret.
	// For receipt generation, we need a secret known to the server.
	// We can use the JWT_SECRET from env + TenantID for signing receipts, ensuring only we can generate them.
	// Or use the tenant's API Secret if we stored it?
	// Security best practice: Derived key. 
	// Let's use a derived key from system secret and tenant ID.
	systemSecret := "system-secret-placeholder" // specific to environment
	tenantSecret := systemSecret + ":" + tenantID.String()

	token, err := h.service.CreateToken(c.Request.Context(), tenantID, tenantSecret, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "create_token_failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, token)
}

// RevokeToken handles POST /v1/consent/tokens/:id/revoke
func (h *Handler) RevokeToken(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "Token ID must be a valid UUID",
		})
		return
	}

	var input RevokeTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	tenantID, _ := middleware.GetTenantID(c)

	token, err := h.service.RevokeToken(c.Request.Context(), tenantID, id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "revoke_token_failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, token)
}

// GetToken handles GET /v1/consent/tokens/:id
func (h *Handler) GetToken(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "Token ID must be a valid UUID",
		})
		return
	}

	tenantID, _ := middleware.GetTenantID(c)

	token, err := h.service.GetToken(c.Request.Context(), tenantID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "not_found",
			"message": "Token not found",
		})
		return
	}

	c.JSON(http.StatusOK, token)
}
