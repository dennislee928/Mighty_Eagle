package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dennislee928/mighty-eagle/api-go/internal/persona"
)

// WorldIDProvider implements verification using World ID
type WorldIDProvider struct {
	AppID      string
	APIKey     string
	HTTPClient *http.Client
}

// NewWorldIDProvider creates a new World ID provider
func NewWorldIDProvider() *WorldIDProvider {
	return &WorldIDProvider{
		AppID:  os.Getenv("WORLDID_APP_ID"),
		APIKey: os.Getenv("WORLDID_API_KEY"),
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Name returns the provider name
func (p *WorldIDProvider) Name() string {
	return "worldid"
}

// Verify verifies a World ID proof
func (p *WorldIDProvider) Verify(ctx context.Context, input persona.VerificationInput) (*persona.VerificationResult, error) {
	// In a real implementation, we would verify the ZK proof here
	// For MVP, we'll verify the proof against the World ID API
	
	// Extract proof from metadata
	proof, ok := input.Metadata["proof"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing proof in metadata")
	}

	payload := map[string]interface{}{
		"merkle_root":    proof["merkle_root"],
		"nullifier_hash": proof["nullifier_hash"],
		"action":         input.Metadata["action"],
		"signal":         input.Metadata["signal"],
		"proof":          proof["proof"],
	}

	appID := p.AppID
	if appID == "" {
		appID = input.Metadata["app_id"].(string) // Allow override or from input
	}
	
	verifyURL := fmt.Sprintf("https://developer.worldcoin.org/api/v1/verify/%s", appID)
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", verifyURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if p.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.APIKey)
	}

	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call World ID API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return &persona.VerificationResult{
			Status:          persona.StatusFailed,
			ConfidenceScore: 0,
			Error:           fmt.Sprintf("World ID API error: %v", errResp),
		}, nil
	}

	// Success
	now := time.Now()
	// World ID verifications are typically good effectively indefinitely for uniqueness, 
	// but applications may want to re-verify occasionally. Setting 1 year for now.
	expiresAt := now.AddDate(1, 0, 0)
	
	return &persona.VerificationResult{
		Status:          persona.StatusVerified,
		ConfidenceScore: 100.0,
		ProviderData:    map[string]interface{}{"verified": true},
		ProofHash:       fmt.Sprintf("%v", proof["nullifier_hash"]),
		VerifiedAt:      &now,
		ExpiresAt:       &expiresAt,
	}, nil
}
