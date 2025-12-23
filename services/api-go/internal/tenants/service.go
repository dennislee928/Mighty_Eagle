package tenants

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/dennislee928/mighty-eagle/api-go/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service provides tenant management functionality
type Service struct {
	db *gorm.DB
}

// NewService creates a new tenant service
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// GenerateAPIKey generates a new API key with prefix
func GenerateAPIKey(tier string) (string, error) {
	// Generate random bytes
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode as hex
	key := hex.EncodeToString(randomBytes)

	// Add prefix based on tier
	var prefix string
	switch tier {
	case "lite":
		prefix = "me_lite_sk_"
	case "pro":
		prefix = "me_pro_sk_"
	case "enterprise":
		prefix = "me_ent_sk_"
	default:
		prefix = "me_sk_"
	}

	return prefix + key, nil
}

// HashAPISecret creates a hash of the API secret for webhook signing
func HashAPISecret(secret string) string {
	hash := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(hash[:])
}

// GenerateAPISecret generates a random secret for webhook signing
func GenerateAPISecret() (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(randomBytes), nil
}

// CreateTenantInput represents data for creating a tenant
type CreateTenantInput struct {
	Name string `json:"name" binding:"required"`
	Tier string `json:"tier" binding:"required,oneof=lite pro enterprise"`
}

// CreateTenant creates a new tenant
func (s *Service) CreateTenant(input CreateTenantInput) (*models.Tenant, error) {
	// Generate API key
	apiKey, err := GenerateAPIKey(input.Tier)
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	// Generate API secret
	apiSecret, err := GenerateAPISecret()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API secret: %w", err)
	}

	// Hash the secret for storage
	apiSecretHash := HashAPISecret(apiSecret)

	tenant := models.Tenant{
		ID:            uuid.New(),
		Name:          input.Name,
		Tier:          input.Tier,
		Status:        "active",
		APIKey:        apiKey,
		APISecretHash: apiSecretHash,
	}

	if err := s.db.Create(&tenant).Error; err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	return &tenant, nil
}

// GetTenant retrieves a tenant by ID
func (s *Service) GetTenant(tenantID uuid.UUID) (*models.Tenant, error) {
	var tenant models.Tenant
	if err := s.db.Where("id = ? AND status = ?", tenantID, "active").First(&tenant).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

// GetTenantByAPIKey retrieves a tenant by API key
func (s *Service) GetTenantByAPIKey(apiKey string) (*models.Tenant, error) {
	var tenant models.Tenant
	if err := s.db.Where("api_key = ? AND status = ?", apiKey, "active").First(&tenant).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

// RegenerateAPIKey regenerates API key for a tenant
func (s *Service) RegenerateAPIKey(tenantID uuid.UUID) (*models.Tenant, error) {
	tenant, err := s.GetTenant(tenantID)
	if err != nil {
		return nil, err
	}

	// Generate new API key
	newAPIKey, err := GenerateAPIKey(tenant.Tier)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new API key: %w", err)
	}

	// Update tenant
	tenant.APIKey = newAPIKey
	if err := s.db.Save(tenant).Error; err != nil {
		return nil, fmt.Errorf("failed to update API key: %w", err)
	}

	return tenant, nil
}

// RotateAPISecret rotates the API secret for webhook signing
func (s *Service) RotateAPISecret(tenantID uuid.UUID) (*models.Tenant, error) {
	tenant, err := s.GetTenant(tenantID)
	if err != nil {
		return nil, err
	}

	// Generate new API secret
	newAPISecret, err := GenerateAPISecret()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new API secret: %w", err)
	}

	// Hash and update
	tenant.APISecretHash = HashAPISecret(newAPISecret)
	if err := s.db.Save(tenant).Error; err != nil {
		return nil, fmt.Errorf("failed to update API secret: %w", err)
	}

	return tenant, nil
}
