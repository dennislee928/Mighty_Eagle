package persona

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dennislee928/mighty-eagle/api-go/internal/audit"
	"github.com/dennislee928/mighty-eagle/api-go/internal/models"
	"github.com/dennislee928/mighty-eagle/api-go/internal/webhooks"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service provides persona verification functionality
type Service struct {
	db        *gorm.DB
	providers map[string]VerificationProvider
	audit     *audit.Logger
	webhooks  *webhooks.Service
}

// NewService creates a new persona service
func NewService(db *gorm.DB, audit *audit.Logger, webhooks *webhooks.Service) *Service {
	return &Service{
		db:        db,
		providers: make(map[string]VerificationProvider),
		audit:     audit,
		webhooks:  webhooks,
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

	// Dispatch Webhook
	go func() {
		// Reconstruct event log metadata for webhook payload
		meta, _ := json.Marshal(map[string]interface{}{
			"provider": input.Provider,
			"score":    result.ConfidenceScore,
			"error":    result.Error,
		})

		dispatchEvt := models.EventLog{
			ID:           uuid.New(), // Temporary ID for dispatch if not persisted yet
			TenantID:     tenantID,
			EventType:    eventType,
			CreatedAt:    time.Now(),
			SubjectID:    &input.SubjectID,
			ResourceID:   &verification.ID,
			Metadata:     string(meta),
		}
		
		if err := s.webhooks.DispatchEvent(context.Background(), dispatchEvt); err != nil {
			fmt.Printf("Failed to dispatch webhook: %v\n", err)
		}
	}()

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
