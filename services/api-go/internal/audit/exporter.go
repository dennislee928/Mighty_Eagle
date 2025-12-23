package audit

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dennislee928/mighty-eagle/api-go/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Exporter manages audit export jobs
type Exporter struct {
	db      *gorm.DB
	logger  *Logger
	billing BillingService
}

type BillingService interface {
	ReportUsage(ctx context.Context, tenantID uuid.UUID, metric string, amount int) error
}

// NewExporter creates a new exporter service
func NewExporter(db *gorm.DB, logger *Logger, billing BillingService) *Exporter {
	return &Exporter{db: db, logger: logger, billing: billing}
}

// CreateExportJob creates a new export job and starts processing
func (e *Exporter) CreateExportJob(ctx context.Context, tenantID uuid.UUID, format string, startDate, endDate time.Time) (*models.AuditExportJob, error) {
	if format != "csv" && format != "json" {
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	job := models.AuditExportJob{
		TenantID:  tenantID,
		Format:    format,
		Status:    "pending",
		StartDate: startDate,
		EndDate:   endDate,
		CreatedAt: time.Now(),
		ExpiresAt: &expiresAt, // Download valid for 24h
		Parameters: string(convertMapToJSON(map[string]interface{}{
			"start_date": startDate,
			"end_date":   endDate,
		})),
	}

	if err := e.db.Create(&job).Error; err != nil {
		return nil, err
	}

	// Start Async Processing
	go e.processJob(job.ID, tenantID, format, startDate, endDate)

	// Report Usage
	go func() {
		if e.billing != nil {
			e.billing.ReportUsage(context.Background(), tenantID, "exports", 1)
		}
	}()

	return &job, nil
}

// GetExportJob retrieves a job status
func (e *Exporter) GetExportJob(ctx context.Context, tenantID uuid.UUID, jobID uuid.UUID) (*models.AuditExportJob, error) {
	var job models.AuditExportJob
	if err := e.db.Where("id = ? AND tenant_id = ?", jobID, tenantID).First(&job).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (e *Exporter) processJob(jobID, tenantID uuid.UUID, format string, startDate, endDate time.Time) {
	// Simple in-memory processing for MVP.
	
	ctx := context.Background()

	// 1. Fetch Logs
	startStr := startDate.Format(time.RFC3339)
	endStr := endDate.Format(time.RFC3339)
	
	logs, _, err := e.logger.GetEventLog(ctx, GetEventLogInput{
		TenantID:  tenantID,
		StartDate: &startStr,
		EndDate:   &endStr,
		Limit:     10000,
	})

	if err != nil {
		e.failJob(jobID, err.Error())
		return
	}

	// 2. Generate Content
	var content []byte
	if format == "csv" {
		content, err = e.generateCSV(logs)
	} else {
		content, err = json.MarshalIndent(logs, "", "  ")
	}

	if err != nil {
		e.failJob(jobID, fmt.Sprintf("failed to generate content: %v", err))
		return
	}

	// 3. Store result as Data URI (MVP Hack)
	dataURI := fmt.Sprintf("data:text/%s;charset=utf-8;base64,%s", format, content) 
	
	e.completeJob(jobID, dataURI, int64(len(content)), len(logs))
}

func (e *Exporter) generateCSV(logs []models.EventLog) ([]byte, error) {
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	
	// Header
	w.Write([]string{"ID", "Time", "Event", "Subject", "Resource", "Metadata"})
	
	for _, l := range logs {
		subj := ""
		if l.SubjectID != nil {
			subj = *l.SubjectID
		}
		res := ""
		if l.ResourceID != nil {
			res = (*l.ResourceID).String()
		}
		
		w.Write([]string{
			l.ID.String(),
			l.CreatedAt.Format(time.RFC3339),
			l.EventType,
			subj,
			res,
			l.Metadata,
		})
	}
	w.Flush()
	return b.Bytes(), nil
}

func (e *Exporter) failJob(id uuid.UUID, reason string) {
	e.db.Model(&models.AuditExportJob{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        "failed",
		"error_message": reason,
	})
}

func (e *Exporter) completeJob(id uuid.UUID, url string, size int64, count int) {
	e.db.Model(&models.AuditExportJob{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":           "completed",
		"download_url":     url,
		"file_size_bytes":  size,
		"record_count":     count,
		"completed_at":     time.Now(),
	})
}

func convertMapToJSON(m map[string]interface{}) []byte {
	b, _ := json.Marshal(m)
	return b
}
