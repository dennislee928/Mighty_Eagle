package models

import (
	"time"

	"github.com/google/uuid"
)

// Tenant represents a platform workspace
type Tenant struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name          string     `gorm:"not null" json:"name"`
	Tier          string     `gorm:"not null;default:'lite'" json:"tier"` // lite, pro, enterprise
	Status        string     `gorm:"not null;default:'active'" json:"status"` // active, suspended, deleted
	APIKey        string     `gorm:"not null;unique" json:"api_key"`
	APISecretHash string     `gorm:"not null" json:"-"` // Hidden from JSON
	CreatedAt     time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

// TableName overrides the table name
func (Tenant) TableName() string {
	return "tenants"
}

// EventLog represents an append-only audit event
type EventLog struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	TenantID     uuid.UUID `gorm:"type:uuid;not null" json:"tenant_id"`
	EventType    string    `gorm:"not null" json:"event_type"`
	EventVersion string    `gorm:"not null;default:'v1'" json:"event_version"`
	ActorID      *string   `json:"actor_id,omitempty"`
	SubjectID    *string   `json:"subject_id,omitempty"`
	ResourceType *string   `json:"resource_type,omitempty"`
	ResourceID   *uuid.UUID `gorm:"type:uuid" json:"resource_id,omitempty"`
	Metadata     string    `gorm:"type:jsonb;not null;default:'{}'" json:"metadata"`
	IPAddress    *string   `gorm:"type:inet" json:"ip_address,omitempty"`
	UserAgent    *string   `json:"user_agent,omitempty"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP;not null" json:"created_at"`
}

// TableName overrides the table name
func (EventLog) TableName() string {
	return "event_log"
}

// PersonaVerification represents a persona verification record
type PersonaVerification struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	TenantID         uuid.UUID  `gorm:"type:uuid;not null" json:"tenant_id"`
	SubjectID        string     `gorm:"not null" json:"subject_id"`
	Provider         string     `gorm:"not null" json:"provider"`
	Status           string     `gorm:"not null;default:'pending'" json:"status"` // pending, verified, failed, expired
	ConfidenceScore  *float64   `gorm:"type:decimal(5,2)" json:"confidence_score,omitempty"`
	VerificationData string     `gorm:"type:jsonb;not null;default:'{}'" json:"verification_data"`
	ProofHash        *string    `json:"proof_hash,omitempty"`
	ExpiresAt        *time.Time `json:"expires_at,omitempty"`
	VerifiedAt       *time.Time `json:"verified_at,omitempty"`
	CreatedAt        time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName overrides the table name
func (PersonaVerification) TableName() string {
	return "persona_verifications"
}

// ConsentToken represents a consent token
type ConsentToken struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	TenantID         uuid.UUID  `gorm:"type:uuid;not null" json:"tenant_id"`
	Parties          string     `gorm:"type:text[];not null" json:"parties"` // Array as JSON string
	Scope            string     `gorm:"not null" json:"scope"`
	TokenHash        string     `gorm:"not null;unique" json:"token_hash"`
	ReceiptSignature string     `gorm:"not null" json:"receipt_signature"`
	Status           string     `gorm:"not null;default:'active'" json:"status"` // active, revoked, expired
	Metadata         string     `gorm:"type:jsonb;default:'{}'" json:"metadata"`
	IssuedAt         time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"issued_at"`
	ExpiresAt        time.Time  `gorm:"not null" json:"expires_at"`
	RevokedAt        *time.Time `json:"revoked_at,omitempty"`
	RevokedBy        *string    `json:"revoked_by,omitempty"`
	RevokeReason     *string    `json:"revoke_reason,omitempty"`
	CreatedAt        time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName overrides the table name
func (ConsentToken) TableName() string {
	return "consent_tokens"
}

// ReputationScore represents a reputation score
type ReputationScore struct {
	ID                    uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	TenantID              uuid.UUID  `gorm:"type:uuid;not null" json:"tenant_id"`
	SubjectID             string     `gorm:"not null" json:"subject_id"`
	Score                 float64    `gorm:"type:decimal(5,2);not null" json:"score"`
	VerifiedStatusWeight  float64    `gorm:"type:decimal(5,2);default:0" json:"verified_status_weight"`
	AccountAgeWeight      float64    `gorm:"type:decimal(5,2);default:0" json:"account_age_weight"`
	DisputeSignalsWeight  float64    `gorm:"type:decimal(5,2);default:0" json:"dispute_signals_weight"`
	Metadata              string     `gorm:"type:jsonb;default:'{}'" json:"metadata"`
	CalculatedAt          time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"calculated_at"`
	ExpiresAt             *time.Time `json:"expires_at,omitempty"`
	CreatedAt             time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt             time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName overrides the table name
func (ReputationScore) TableName() string {
	return "reputation_scores"
}

// WebhookEndpoint represents a webhook endpoint configuration
type WebhookEndpoint struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null" json:"tenant_id"`
	URL       string    `gorm:"not null" json:"url"`
	Secret    string    `gorm:"not null" json:"-"` // Hidden from JSON
	Enabled   bool      `gorm:"default:true" json:"enabled"`
	Events    string    `gorm:"type:text[];not null" json:"events"` // Array as JSON string
	Metadata  string    `gorm:"type:jsonb;default:'{}'" json:"metadata"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName overrides the table name
func (WebhookEndpoint) TableName() string {
	return "webhook_endpoints"
}

// WebhookDelivery represents a webhook delivery attempt
type WebhookDelivery struct {
	ID                uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	WebhookEndpointID uuid.UUID  `gorm:"type:uuid;not null" json:"webhook_endpoint_id"`
	EventID           uuid.UUID  `gorm:"type:uuid;not null" json:"event_id"`
	Status            string     `gorm:"not null;default:'pending'" json:"status"` // pending, success, failed, dlq
	AttemptCount      int        `gorm:"default:0" json:"attempt_count"`
	MaxAttempts       int        `gorm:"default:3" json:"max_attempts"`
	RequestPayload    string     `gorm:"type:jsonb;not null" json:"request_payload"`
	RequestHeaders    *string    `gorm:"type:jsonb" json:"request_headers,omitempty"`
	ResponseStatus    *int       `json:"response_status,omitempty"`
	ResponseBody      *string    `json:"response_body,omitempty"`
	ErrorMessage      *string    `json:"error_message,omitempty"`
	NextRetryAt       *time.Time `json:"next_retry_at,omitempty"`
	CompletedAt       *time.Time `json:"completed_at,omitempty"`
	CreatedAt         time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName overrides the table name
func (WebhookDelivery) TableName() string {
	return "webhook_deliveries"
}
