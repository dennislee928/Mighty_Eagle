package persona

import (
	"context"
	"time"
)

// VerificationStatus represents the status of a verification
type VerificationStatus string

const (
	StatusPending  VerificationStatus = "pending"
	StatusVerified VerificationStatus = "verified"
	StatusFailed   VerificationStatus = "failed"
	StatusExpired  VerificationStatus = "expired"
)

// VerificationInput represents the input required for verification
type VerificationInput struct {
	SubjectID string                 `json:"subject_id"`
	Provider  string                 `json:"provider"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// VerificationResult represents the result of a verification attempt
type VerificationResult struct {
	Status          VerificationStatus     `json:"status"`
	ConfidenceScore float64                `json:"confidence_score"`
	ProviderData    map[string]interface{} `json:"provider_data"`
	ProofHash       string                 `json:"proof_hash,omitempty"`
	VerifiedAt      *time.Time             `json:"verified_at,omitempty"`
	ExpiresAt       *time.Time             `json:"expires_at,omitempty"`
	Error           string                 `json:"error,omitempty"`
}

// VerificationProvider defines the interface for persona verification providers
type VerificationProvider interface {
	// Verify initiates or checks a verification request
	Verify(ctx context.Context, input VerificationInput) (*VerificationResult, error)
	
	// Name returns the provider name
	Name() string
}
