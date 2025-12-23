package persona

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dennislee928/mighty-eagle/api-go/internal/audit"
	"github.com/dennislee928/mighty-eagle/api-go/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service provides persona verification functionality
type Service struct {
	db        *gorm.DB
	providers map[string]VerificationProvider
	audit     *audit.Logger
}

// NewService creates a new persona service
func NewService(db *gorm.DB, audit *audit.Logger) *Service {
	return &Service{
		db:        db,
		providers: make(map[string]VerificationProvider),
		audit:     audit,
	}
}

// RegisterProvider registers a verification provider
func (s *Service) RegisterProvider(provider VerificationProvider) {
	s.providers[provider.Name()] = provider
}

// CreateVerification initiates a verification request
func (s *Service) CreateVerification(ctx context.Context, tenantID uuid.UUID, input VerificationInput) (*models.PersonaVerification, error) {
	// Find provider
	provider, ok := s.providers[input.Provider]
	if !ok {
		return nil, fmt.Errorf("provider '%s' not supported", input.Provider)
	}

	// Perform verification
	result, err := provider.Verify(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("verification failed: %w", err)
	}

	// Marshal provider data
	providerDataJSON, err := json.Marshal(result.ProviderData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal provider data: %w", err)
	}

	// Create record
	verification := models.PersonaVerification{
		TenantID:         tenantID,
		SubjectID:        input.SubjectID,
		Provider:         input.Provider,
		Status:           string(result.Status),
		ConfidenceScore:  &result.ConfidenceScore,
		VerificationData: string(providerDataJSON),
		ProofHash:        &result.ProofHash,
		ExpiresAt:        result.ExpiresAt,
		VerifiedAt:       result.VerifiedAt,
	}

	if err := s.db.Create(&verification).Error; err != nil {
		return nil, fmt.Errorf("failed to create verification record: %w", err)
	}

	// Log audit event
	eventType := "persona.verified"
	if result.Status == StatusFailed {
		eventType = "persona.failed"
	} else if result.Status == StatusPending {
		eventType = "persona.pending"
	}

	s.audit.LogEvent(ctx, audit.LogEventInput{
		TenantID:     tenantID,
		EventType:    eventType,
		SubjectID:    &input.SubjectID,
		ResourceType: stringPtr("verification"),
		ResourceID:   &verification.ID,
		Metadata: map[string]interface{}{
			"provider": input.Provider,
			"score":    result.ConfidenceScore,
			"error":    result.Error,
		},
	})

	// TODO: Trigger webhook (M4)

	return &verification, nil
}

// GetVerification retrieves a verification by ID
func (s *Service) GetVerification(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*models.PersonaVerification, error) {
	var verification models.PersonaVerification
	if err := s.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&verification).Error; err != nil {
		return nil, err
	}
	return &verification, nil
}

func stringPtr(s string) *string {
	return &s
}
