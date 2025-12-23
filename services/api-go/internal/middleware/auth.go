package middleware

import (
	"net/http"

	"github.com/dennislee928/mighty-eagle/api-go/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthMiddleware validates API key and injects tenant context
func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract API key from header
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing_api_key",
				"message": "API key is required. Include 'X-API-Key' header.",
			})
			c.Abort()
			return
		}

		// Validate API key and retrieve tenant
		var tenant models.Tenant
		result := db.Where("api_key = ? AND status = ?", apiKey, "active").First(&tenant)
		
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "invalid_api_key",
					"message": "Invalid or inactive API key.",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "auth_error",
					"message": "Failed to authenticate request.",
				})
			}
			c.Abort()
			return
		}

		// Inject tenant into context
		c.Set("tenant_id", tenant.ID)
		c.Set("tenant", tenant)
		c.Set("tenant_tier", tenant.Tier)

		c.Next()
	}
}

// GetTenantID retrieves tenant ID from context
func GetTenantID(c *gin.Context) (uuid.UUID, bool) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		return uuid.Nil, false
	}
	id, ok := tenantID.(uuid.UUID)
	return id, ok
}

// GetTenant retrieves full tenant from context
func GetTenant(c *gin.Context) (*models.Tenant, bool) {
	tenant, exists := c.Get("tenant")
	if !exists {
		return nil, false
	}
	t, ok := tenant.(models.Tenant)
	if !ok {
		return nil, false
	}
	return &t, true
}

// OptionalAuthMiddleware validates API key if present, otherwise continues
func OptionalAuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.Next()
			return
		}

		// If API key is provided, validate it
		var tenant models.Tenant
		result := db.Where("api_key = ? AND status = ?", apiKey, "active").First(&tenant)
		
		if result.Error == nil {
			c.Set("tenant_id", tenant.ID)
			c.Set("tenant", tenant)
			c.Set("tenant_tier", tenant.Tier)
		}

		c.Next()
	}
}
