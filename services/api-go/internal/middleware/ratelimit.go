package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int
	BurstSize         int
}

// TierLimits defines rate limits for each tier
var TierLimits = map[string]RateLimitConfig{
	"lite":       {RequestsPerMinute: 60, BurstSize: 10},
	"pro":        {RequestsPerMinute: 600, BurstSize: 50},
	"enterprise": {RequestsPerMinute: 6000, BurstSize: 200},
}

// RateLimitMiddleware implements token bucket rate limiting using Redis
func RateLimitMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant tier from context (set by AuthMiddleware)
		tier, exists := c.Get("tenant_tier")
		if !exists {
			tier = "lite" // Default to lite tier if not authenticated
		}

		tierStr, ok := tier.(string)
		if !ok {
			tierStr = "lite"
		}

		config, ok := TierLimits[tierStr]
		if !ok {
			config = TierLimits["lite"]
		}

		// Get tenant ID for rate limit key
		tenantID, _ := c.Get("tenant_id")
		key := fmt.Sprintf("ratelimit:%v:%s", tenantID, time.Now().Format("2006-01-02-15-04"))

		ctx := context.Background()

		// Increment counter
		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			// If Redis is down, allow the request but log the error
			c.Next()
			return
		}

		// Set expiry on first increment
		if count == 1 {
			redisClient.Expire(ctx, key, time.Minute)
		}

		// Check if limit exceeded
		if count > int64(config.RequestsPerMinute) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate_limit_exceeded",
				"message": fmt.Sprintf("Rate limit exceeded. Maximum %d requests per minute for %s tier.", config.RequestsPerMinute, tierStr),
				"retry_after": 60,
			})
			c.Abort()
			return
		}

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", config.RequestsPerMinute))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", config.RequestsPerMinute-int(count)))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Minute).Unix()))

		c.Next()
	}
}

// GlobalRateLimitMiddleware applies a global rate limit per IP
func GlobalRateLimitMiddleware(redisClient *redis.Client, requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("global_ratelimit:%s:%s", ip, time.Now().Format("2006-01-02-15-04"))

		ctx := context.Background()

		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			redisClient.Expire(ctx, key, time.Minute)
		}

		if count > int64(requestsPerMinute) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "global_rate_limit_exceeded",
				"message": "Too many requests. Please try again later.",
				"retry_after": 60,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
