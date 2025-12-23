package billing

import (
	"net/http"

	"github.com/dennislee928/mighty-eagle/api-go/internal/middleware"
	"github.com/gin-gonic/gin"
)

// CheckEntitlementMiddleware returns a middleware that checks if the tenant has quota for a metric
func CheckEntitlementMiddleware(service *Service, metric string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenant, ok := middleware.GetTenant(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		allowed, err := service.CheckEntitlement(c.Request.Context(), *tenant, metric)
		if err != nil {
			// Log error but maybe fail open or closed? Fail closed (deny) is safer for billing.
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "entitlement_check_failed"})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusPaymentRequired, gin.H{
				"error":   "quota_exceeded",
				"message": "You have exceeded your usage limit for " + metric + ". Please upgrade your plan.",
			})
			return
		}

		c.Next()
	}
}
