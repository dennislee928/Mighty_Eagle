package audit

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dennislee928/mighty-eagle/api-go/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Logger provides event logging functionality
type Logger struct {
	db *gorm.DB
}

// NewLogger creates a new audit logger
func NewLogger(db *gorm.DB) *Logger {
	return &Logger{db: db}
}

// LogEventInput represents data for logging an event
type LogEventInput struct {
	TenantID     uuid.UUID
	EventType    string
	ActorID      *string
	SubjectID    *string
	ResourceType *string
	ResourceID   *uuid.UUID
	Metadata     map[string]interface{}
	IPAddress    *string
	UserAgent    *string
}

// LogEvent logs an event to the audit trail
func (l *Logger) LogEvent(ctx context.Context, input LogEventInput) error {
	// Convert metadata to JSON
	metadataJSON, err := json.Marshal(input.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	event := models.EventLog{
		TenantID:     input.TenantID,
		EventType:    input.EventType,
		EventVersion: "v1",
		ActorID:      input.ActorID,
		SubjectID:    input.SubjectID,
		ResourceType: input.ResourceType,
		ResourceID:   input.ResourceID,
		Metadata:     string(metadataJSON),
		IPAddress:    input.IPAddress,
		UserAgent:    input.UserAgent,
	}

	if err := l.db.WithContext(ctx).Create(&event).Error; err != nil {
		return fmt.Errorf("failed to create event log: %w", err)
	}

	return nil
}

// LogEventFromContext is a convenience method that extracts context data from Gin
func (l *Logger) LogEventFromContext(c *gin.Context, input LogEventInput) error {
	// Extract tenant ID if not provided
	if input.TenantID == uuid.Nil {
		tenantID, exists := c.Get("tenant_id")
		if exists {
			if tid, ok := tenantID.(uuid.UUID); ok {
				input.TenantID = tid
			}
		}
	}

	// Extract IP and User-Agent if not provided
	if input.IPAddress == nil {
		ip := c.ClientIP()
		input.IPAddress = &ip
	}

	if input.UserAgent == nil {
		ua := c.Request.UserAgent()
		input.UserAgent = &ua
	}

	return l.LogEvent(c.Request.Context(), input)
}

// GetEventLog retrieves events with pagination and filtering
type GetEventLogInput struct {
	TenantID   uuid.UUID
	EventTypes []string
	StartDate  *string
	EndDate    *string
	Limit      int
	Offset     int
}

// GetEventLog retrieves event logs
func (l *Logger) GetEventLog(ctx context.Context, input GetEventLogInput) ([]models.EventLog, int64, error) {
	query := l.db.WithContext(ctx).Where("tenant_id = ?", input.TenantID)

	// Apply filters
	if len(input.EventTypes) > 0 {
		query = query.Where("event_type IN ?", input.EventTypes)
	}

	if input.StartDate != nil {
		query = query.Where("created_at >= ?", *input.StartDate)
	}

	if input.EndDate != nil {
		query = query.Where("created_at <= ?", *input.EndDate)
	}

	// Get total count
	var total int64
	if err := query.Model(&models.EventLog{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count events: %w", err)
	}

	// Apply pagination
	if input.Limit == 0 {
		input.Limit = 100
	}
	query = query.Limit(input.Limit).Offset(input.Offset)

	// Order by created_at DESC
	query = query.Order("created_at DESC")

	// Execute query
	var events []models.EventLog
	if err := query.Find(&events).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve events: %w", err)
	}

	return events, total, nil
}
