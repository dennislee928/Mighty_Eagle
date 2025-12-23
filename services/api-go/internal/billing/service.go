package billing

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dennislee928/mighty-eagle/api-go/internal/audit"
	"github.com/dennislee928/mighty-eagle/api-go/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Tier limits
var (
	Limits = map[string]map[string]int64{
		"lite": {
			"verifications": 100,
			"exports":       10,
		},
		"pro": {
			"verifications": 10000,
			"exports":       100,
		},
		"enterprise": {
			"verifications": 1000000,
			"exports":       10000,
		},
	}
)

// Service manages billing and entitlements
type Service struct {
	db          *gorm.DB
	redisClient *redis.Client
	audit       *audit.Logger
}

// NewService creates a new billing service
func NewService(db *gorm.DB, redisClient *redis.Client, audit *audit.Logger) *Service {
	return &Service{db: db, redisClient: redisClient, audit: audit}
}

// ReportUsage increments usage for a metric
func (s *Service) ReportUsage(ctx context.Context, tenantID uuid.UUID, metric string, amount int) error {
	// For MVP: Increment Redis counter for current month.
	// Format: usage:{tenant}:{metric}:YYYY-MM
	currentMonth := time.Now().Format("2006-01")
	key := fmt.Sprintf("usage:%s:%s:%s", tenantID, metric, currentMonth)

	pipe := s.redisClient.Pipeline()
	pipe.IncrBy(ctx, key, int64(amount))
	pipe.Expire(ctx, key, 60*24*time.Hour) // Keep for 60 days
	_, err := pipe.Exec(ctx)

	// Async: Persist to DB periodically? 
	// For M5 MVP: We rely on Redis for real-time checks. 
	// We should also update `usage_metrics` table? 
	// The schema has `usage_metrics` table. Let's update it or Insert new row.
	// Since report is frequent (per request), DB write is expensive.
	// We'll skip DB write on every request for MVP and trust Redis + Async sync (not implemented yet).
	// OR: Just write to DB if metric is critical (like export). Verifications are high volume.
	
	if err != nil {
		fmt.Printf("Failed to report usage: %v\n", err)
	}
	return err
}

// CheckEntitlement checks if a tenant can perform an action based on limits
func (s *Service) CheckEntitlement(ctx context.Context, tenant models.Tenant, metric string) (bool, error) {
	tier := strings.ToLower(tenant.Tier)
	
	// Default to lite if unknown
	limitMap, ok := Limits[tier]
	if !ok {
		limitMap = Limits["lite"]
	}

	limit, ok := limitMap[metric]
	if !ok {
		// No limit defined for metric, allow
		return true, nil
	}
	
	if limit == -1 { // Unlimited convention
		return true, nil
	}

	// Check usage
	currentMonth := time.Now().Format("2006-01")
	key := fmt.Sprintf("usage:%s:%s:%s", tenant.ID, metric, currentMonth)
	
	val, err := s.redisClient.Get(ctx, key).Int64()
	if err == redis.Nil {
		val = 0
	} else if err != nil {
		return false, fmt.Errorf("failed to check usage: %w", err)
	}

	if val >= limit {
		return false, nil
	}

	return true, nil
}

// GetCurrentUsage retrieves usage stats for the tenant
func (s *Service) GetCurrentUsage(ctx context.Context, tenant models.Tenant) (map[string]map[string]interface{}, error) {
	currentMonth := time.Now().Format("2006-01")
	metrics := []string{"verifications", "exports"}
	
	result := make(map[string]map[string]interface{})
	tier := strings.ToLower(tenant.Tier)
	limitMap := Limits[tier]
	if limitMap == nil {
		limitMap = Limits["lite"]
	}

	for _, m := range metrics {
		key := fmt.Sprintf("usage:%s:%s:%s", tenant.ID, m, currentMonth)
		val, _ := s.redisClient.Get(ctx, key).Int64()
		
		limit := limitMap[m]
		
		result[m] = map[string]interface{}{
			"used":  val,
			"limit": limit,
		}
	}

	return result, nil
}
