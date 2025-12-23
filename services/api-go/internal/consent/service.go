package consent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dennislee928/mighty-eagle/api-go/internal/audit"
	"github.com/dennislee928/mighty-eagle/api-go/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service manages consent tokens
type Service struct {
	db    *gorm.DB
	audit *audit.Logger
}

// NewService creates a new consent service
func NewService(db *gorm.DB, audit *audit.Logger) *Service {
	return &Service{db: db, audit: audit}
}

// CreateTokenInput represents input for creating a consent token
type CreateTokenInput struct {
	Parties   []string               `json:"parties" binding:"required,min=2"`
	Scope     string                 `json:"scope" binding:"required"`
	ExpiresAt time.Time              `json:"expires_at" binding:"required"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// CreateToken issues a new consent token
func (s *Service) CreateToken(ctx context.Context, tenantID uuid.UUID, tenantSecret string, input CreateTokenInput) (*models.ConsentToken, error) {
	// 1. Generate unique hash to prevent duplicates if business rule requires unique active consent per scope
	tokenHash := GenerateTokenHash(tenantID, input.Parties, input.Scope)
	
	// Check for existing active token?
	// For MVP, we'll allow multiple but maybe warn? Or the DB constraint handles uniqueness on token_hash?
	// The schema defines token_hash as unique. So we must ensure it's unique or we fail.
	// We might want to include expiry in hash to allow renewal?
	// Let's assume (Tenant + Parties + Scope) must be unique. If re-issuing, revoke old one first?
	// Or add randomness/timestamp to hash to allow multiple?
	// Decision: Add timestamp to hash to allow multiple consents (history) but receipt creates binding.
	// Wait, schema says `token_hash VARCHAR(255) NOT NULL UNIQUE`.
	// Let's add timestamp to the hash input to make it unique per issuance.
	tokenHashData := fmt.Sprintf("%s:%v:%s:%d", tenantID, input.Parties, input.Scope, time.Now().UnixNano())
	tokenHash = GenerateTokenHash(tenantID, []string{tokenHashData}, "")

	// 2. Create Receipt
	tokenID := uuid.New()
	receiptPayload := ReceiptPayload{
		TokenID:   tokenID,
		Parties:   input.Parties,
		Scope:     input.Scope,
		ExpiresAt: input.ExpiresAt,
		TenantID:  tenantID,
	}
	
	receipt, err := GenerateReceipt(receiptPayload, tenantSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate receipt: %w", err)
	}

	// 3. Prepare data
	// Format the string manually as Postgres array format: "{a,b}"
	partiesPGStr := toPostgresArray(input.Parties)

	metadataJSON, _ := json.Marshal(input.Metadata)

	token := models.ConsentToken{
		ID:               tokenID,
		TenantID:         tenantID,
		Parties:          partiesPGStr,
		Scope:            input.Scope,
		TokenHash:        tokenHash,
		ReceiptSignature: receipt,
		Status:           "active",
		Metadata:         string(metadataJSON),
		ExpiresAt:        input.ExpiresAt,
		IssuedAt:         time.Now(),
	}

	if err := s.db.Create(&token).Error; err != nil {
		return nil, fmt.Errorf("failed to create consent token: %w", err)
	}

	// 4. Audit Log
	s.audit.LogEvent(ctx, audit.LogEventInput{
		TenantID:     tenantID,
		EventType:    "consent.issued",
		ResourceType: stringPtr("consent_token"),
		ResourceID:   &token.ID,
		Metadata: map[string]interface{}{
			"scope":   input.Scope,
			"parties": input.Parties,
		},
	})

	return &token, nil
}

// RevokeTokenInput represents input for revoking a token
type RevokeTokenInput struct {
	RevokedBy string `json:"revoked_by" binding:"required"`
	Reason    string `json:"reason"`
}

// RevokeToken revokes a consent token
func (s *Service) RevokeToken(ctx context.Context, tenantID uuid.UUID, tokenID uuid.UUID, input RevokeTokenInput) (*models.ConsentToken, error) {
	var token models.ConsentToken
	if err := s.db.Where("id = ? AND tenant_id = ?", tokenID, tenantID).First(&token).Error; err != nil {
		return nil, err
	}

	if token.Status != "active" {
		return nil, fmt.Errorf("token is already %s", token.Status)
	}

	now := time.Now()
	token.Status = "revoked"
	token.RevokedAt = &now
	token.RevokedBy = &input.RevokedBy
	token.RevokeReason = &input.Reason

	if err := s.db.Save(&token).Error; err != nil {
		return nil, fmt.Errorf("failed to revoke token: %w", err)
	}

	// Audit Log
	s.audit.LogEvent(ctx, audit.LogEventInput{
		TenantID:     tenantID,
		EventType:    "consent.revoked",
		ResourceType: stringPtr("consent_token"),
		ResourceID:   &token.ID,
		Metadata: map[string]interface{}{
			"revoked_by": input.RevokedBy,
			"reason":     input.Reason,
		},
	})

	return &token, nil
}

// GetToken retrieves a token
func (s *Service) GetToken(ctx context.Context, tenantID uuid.UUID, tokenID uuid.UUID) (*models.ConsentToken, error) {
	var token models.ConsentToken
	if err := s.db.Where("id = ? AND tenant_id = ?", tokenID, tenantID).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

// Helper to convert slice to Postgres array string format "{a,b}"
// simple implementation, assumes no special chars (comma, brace) in values for MVP
func toPostgresArray(arr []string) string {
	if len(arr) == 0 {
		return "{}"
	}
	// Better implementation would escaping, but for IDs usually safe
	res := "{"
	for i, s := range arr {
		if i > 0 {
			res += ","
		}
		res += fmt.Sprintf("\"%s\"", s)
	}
	res += "}"
	return res
}

func stringPtr(s string) *string {
	return &s
}
