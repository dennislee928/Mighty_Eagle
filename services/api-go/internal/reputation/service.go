package reputation

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dennislee928/mighty-eagle/api-go/internal/audit"
	"github.com/dennislee928/mighty-eagle/api-go/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Service manages reputation scores
type Service struct {
	db          *gorm.DB
	redisClient *redis.Client
	audit       *audit.Logger
	scorer      *Scorer
}

// NewService creates a new reputation service
func NewService(db *gorm.DB, redisClient *redis.Client, audit *audit.Logger) *Service {
	return &Service{
		db:          db,
		redisClient: redisClient,
		audit:       audit,
		scorer:      NewScorer(),
	}
}

// ReputationResult represents the API response
type ReputationResult struct {
	SubjectID      string          `json:"subject_id"`
	Score          float64         `json:"score"`
	Level          string          `json:"level"` // Low, Medium, High, Very High
	Components     ScoreComponents `json:"components"`
	LastCalculated time.Time       `json:"last_calculated"`
	Cached         bool            `json:"cached"`
}

// GetReputation retrieves or calculates the reputation score for a subject
func (s *Service) GetReputation(ctx context.Context, tenantID uuid.UUID, subjectID string) (*ReputationResult, error) {
	cacheKey := fmt.Sprintf("reputation:%s:%s", tenantID, subjectID)

	// 1. Check Cache
	val, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var result ReputationResult
		if err := json.Unmarshal([]byte(val), &result); err == nil {
			result.Cached = true
			return &result, nil
		}
	}

	// 2. Fetch Data
	// Fetch Verifications
	var verifications []models.PersonaVerification
	if err := s.db.Where("tenant_id = ? AND subject_id = ?", tenantID, subjectID).Find(&verifications).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch verification history: %w", err)
	}

	// 3. Calculate Score
	score, components := s.scorer.CalculateScore(verifications, nil)

	// Determine Level
	level := "Low"
	if score >= 80 {
		level = "Very High"
	} else if score >= 60 {
		level = "High"
	} else if score >= 40 {
		level = "Medium"
	}

	result := ReputationResult{
		SubjectID:      subjectID,
		Score:          score,
		Level:          level,
		Components:     components,
		LastCalculated: time.Now(),
		Cached:         false,
	}

	// 4. Cache Result (1 Hour TTL)
	jsonBytes, _ := json.Marshal(result)
	s.redisClient.Set(ctx, cacheKey, jsonBytes, 1*time.Hour)

	// 5. Persist Score Snapshot (optional but good for history)
	// For M3, let's just create/update the reputation_scores table
	// Upsert logic
	snapshot := models.ReputationScore{
		TenantID:  tenantID,
		SubjectID: subjectID,
		Score:     score,
		Level:     level,
		Factors:   string(convertMapToJSON(map[string]interface{}{"components": components})),
		UpdatedAt: time.Now(),
	}
	
	// Assuming unique constraint on (tenant_id, subject_id) in schema logic.
	// Although the schema migration might not have unique constraint explicitly named?
	// Let's replace or create.
	// Actually, the initial schema.sql creates the table but might not enforce uniqueness on current entry.
	// But it's better to log a new entry if we want history. 
	// The schema says `reputation_scores` has an ID.
	// Let's just create a new record for "current state" (or overwrite if we track "current").
	// Implementation choice: Keep track of "latest" in a `current_reputation` or just append to `reputation_scores` as history log?
	// Reviewing schema.sql: "CREATE TABLE reputation_scores" looks like a log with timestamp.
	// So we append.
	
	// Wait, constant writes on every GET might be heavy if un-cached. 
	// But we cache for 1 hour. So writes happen max once per hour per active user.
	// Acceptable for MVP.
	snapshot.ID = uuid.New()
	s.db.Create(&snapshot)

	return &result, nil
}

func convertMapToJSON(m map[string]interface{}) []byte {
	b, _ := json.Marshal(m)
	return b
}
