package consent

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ReceiptPayload represents the data included in a consent receipt
type ReceiptPayload struct {
	TokenID   uuid.UUID `json:"token_id"`
	Parties   []string  `json:"parties"`
	Scope     string    `json:"scope"`
	ExpiresAt time.Time `json:"expires_at"`
	TenantID  uuid.UUID `json:"tenant_id"`
}

// GenerateReceipt generates a signed receipt string
// Format: "v1.{payload_base64}.{signature_hex}"
func GenerateReceipt(payload ReceiptPayload, secret string) (string, error) {
	// 1. Marshal payload
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal receipt payload: %w", err)
	}
	
	// 2. Encode payload as hex (simpler/safer for URL usage than base64 sometimes)
	// Actually, let's use hex for consistency with other parts for now.
	payloadHex := hex.EncodeToString(jsonBytes)
	
	// 3. Sign payload
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payloadHex))
	signature := hex.EncodeToString(h.Sum(nil))
	
	// 4. Construct receipt string
	return fmt.Sprintf("v1.%s.%s", payloadHex, signature), nil
}

// GenerateTokenHash creates a deterministic hash of the consent parameters
// This ensures uniqueness for (Parties + Scope + Tenant) if desired, 
// or allows lookup by hash.
func GenerateTokenHash(tenantID uuid.UUID, parties []string, scope string) string {
	// Sort parties to ensure deterministic ordering? 
	// For MVP, assume caller provides consistent order or we concat simply.
	// Let's just hash the inputs.
	
	data := fmt.Sprintf("%s:%v:%s", tenantID.String(), parties, scope)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// VerifyReceipt verifies a receipt string against a secret
func VerifyReceipt(receipt string, secret string) (bool, *ReceiptPayload, error) {
	// TODO: Parse and verify if needed later
	return false, nil, fmt.Errorf("not implemented")
}
