package router

import (
	"net/http"

	"github.com/dennislee928/mighty-eagle/api-go/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// SetupRouter initializes all routes
func SetupRouter(db *gorm.DB, redisClient *redis.Client) *gin.Engine {
	r := gin.Default()

	// Global middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.GlobalRateLimitMiddleware(redisClient, 1000)) // 1000 requests per minute per IP

	// Health check endpoint (no auth required)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": "mighty-eagle-api",
			"version": "0.1.0",
		})
	})

	// API v1 routes
	v1 := r.Group("/v1")
	{
		// Apply authentication and rate limiting to all v1 routes
		v1.Use(middleware.AuthMiddleware(db))
		v1.Use(middleware.RateLimitMiddleware(redisClient))

		// TODO: Add route groups for each module
		// - Persona verification
		// - Consent tokens
		// - Reputation
		// - Webhooks
		// - Audit exports
		// - Billing

		// Placeholder: API info endpoint
		v1.GET("/info", func(c *gin.Context) {
			tenant, _ := middleware.GetTenant(c)
			c.JSON(http.StatusOK, gin.H{
				"message": "Mighty Eagle Trust Layer API",
				"version": "v1",
				"tenant": gin.H{
					"id":   tenant.ID,
					"name": tenant.Name,
					"tier": tenant.Tier,
				},
			})
		})
	}

	// 404 handler
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "not_found",
			"message": "The requested endpoint does not exist.",
		})
	})

	return r
}
