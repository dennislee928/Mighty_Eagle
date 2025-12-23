package router

import (
	"context"
	"net/http"
	"os"

	"github.com/dennislee928/mighty-eagle/api-go/internal/audit"
	"github.com/dennislee928/mighty-eagle/api-go/internal/consent"
	"github.com/dennislee928/mighty-eagle/api-go/internal/middleware"
	"github.com/dennislee928/mighty-eagle/api-go/internal/persona"
	"github.com/dennislee928/mighty-eagle/api-go/internal/persona/providers"
	"github.com/dennislee928/mighty-eagle/api-go/internal/webhooks"
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

	// Initialize Services
	auditLogger := audit.NewLogger(db)
	webhookService := webhooks.NewService(db)
	
	// Start webhook worker
	go webhookService.Worker(context.Background())
	
	personaService := persona.NewService(db, auditLogger, webhookService)
	personaService.RegisterProvider(providers.NewMockProvider())
	if os.Getenv("WORLDID_APP_ID") != "" {
		personaService.RegisterProvider(providers.NewWorldIDProvider())
	}
	personaHandler := persona.NewHandler(personaService)

	consentService := consent.NewService(db, auditLogger)
	consentHandler := consent.NewHandler(consentService)

	// API v1 routes
	v1 := r.Group("/v1")
	{
		// Apply authentication and rate limiting to all v1 routes
		v1.Use(middleware.AuthMiddleware(db))
		v1.Use(middleware.RateLimitMiddleware(redisClient))

		// Persona verification routes
		v1.POST("/persona/verifications", personaHandler.CreateVerification)
		v1.GET("/persona/verifications/:id", personaHandler.GetVerification)
		
		// Consent token routes
		v1.POST("/consent/tokens", consentHandler.CreateToken)
		v1.POST("/consent/tokens/:id/revoke", consentHandler.RevokeToken)
		v1.GET("/consent/tokens/:id", consentHandler.GetToken)

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
