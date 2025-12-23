package providers

import (
	"context"
	"time"

	"github.com/dennislee928/mighty-eagle/api-go/internal/persona"
)

// MockProvider implements the VerificationProvider interface for testing
type MockProvider struct{}

// NewMockProvider creates a new instance of MockProvider
func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

// Name returns the provider name
func (p *MockProvider) Name() string {
	return "mock"
}

// Verify simulates a verification process
func (p *MockProvider) Verify(ctx context.Context, input persona.VerificationInput) (*persona.VerificationResult, error) {
	// Simulate success/failure based on subject_id pattern
	status := persona.StatusVerified
	confidence := 100.0
	var errorMsg string

	if input.SubjectID == "fail_me" {
		status = persona.StatusFailed
		confidence = 0.0
		errorMsg = "simulation_forced_failure"
	}
	
	now := time.Now()
	// Default expiry 1 year
	expiresAt := now.AddDate(1, 0, 0)
	
	if status == persona.StatusVerified {
		return &persona.VerificationResult{
			Status:          status,
			ConfidenceScore: confidence,
			ProviderData: map[string]interface{}{
				"mock_session_id": "mock_abc123",
				"simulated":       true,
			},
			ProofHash:  "mock_proof_hash_xyz",
			VerifiedAt: &now,
			ExpiresAt:  &expiresAt,
		}, nil
	}

	return &persona.VerificationResult{
		Status:          status,
		ConfidenceScore: confidence,
		Error:           errorMsg,
	}, nil
}
