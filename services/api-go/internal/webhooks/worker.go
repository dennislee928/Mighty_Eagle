package webhooks

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/dennislee928/mighty-eagle/api-go/internal/models"
)

// DispatchEvent creates delivery records for an event
func (s *Service) DispatchEvent(ctx context.Context, eventLog models.EventLog) error {
	// Find matching subscriptions
	var endpoints []models.WebhookEndpoint
	// Postgers array overlaps operator &&
	if err := s.db.Where("tenant_id = ? AND enabled = ? AND ? = ANY(events)", eventLog.TenantID, true, eventLog.EventType).Find(&endpoints).Error; err != nil {
		return err
	}

	if len(endpoints) == 0 {
		return nil
	}

	// Create payload
	payload := map[string]interface{}{
		"id":         eventLog.ID,
		"event":      eventLog.EventType,
		"created_at": eventLog.CreatedAt,
		"data":       json.RawMessage(eventLog.Metadata),
		"resource":   eventLog.ResourceID,
	}
	
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	payloadStr := string(payloadBytes)

	// Create deliveries
	for _, endpoint := range endpoints {
		delivery := models.WebhookDelivery{
			WebhookEndpointID: endpoint.ID,
			EventID:           eventLog.ID,
			Status:            "pending",
			RequestPayload:    payloadStr,
			MaxAttempts:       3,
			NextRetryAt:       &eventLog.CreatedAt, // Process immediately
		}
		
		if err := s.db.Create(&delivery).Error; err != nil {
			log.Printf("Failed to create delivery for endpoint %s: %v", endpoint.ID, err)
			continue
		}
	}

	return nil
}

// Worker processes pending webhook deliveries
func (s *Service) Worker(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.processPendingDeliveries(ctx)
		}
	}
}

func (s *Service) processPendingDeliveries(ctx context.Context) {
	var deliveries []models.WebhookDelivery
	
	// Find pending deliveries due for retry
	if err := s.db.Where("status = ? AND next_retry_at <= ?", "pending", time.Now()).Limit(50).Find(&deliveries).Error; err != nil {
		log.Printf("Error fetching pending deliveries: %v", err)
		return
	}

	for _, d := range deliveries {
		go func(delivery models.WebhookDelivery) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			
			if err := s.Deliver(ctx, delivery.ID); err != nil {
				// Error logging handled in Deliver/recordFailure usually, but log here too
				// log.Printf("Delivery %s failed: %v", delivery.ID, err)
			}
		}(d)
	}
}
