package webhooks

import (
	"net/http"

	"github.com/dennislee928/mighty-eagle/api-go/internal/middleware"
	"github.com/dennislee928/mighty-eagle/api-go/internal/models"
	"github.com/gin-gonic/gin"
)

// Handler manages webhook-related HTTP endpoints
type Handler struct {
	service *Service
}

// NewHandler creates a new webhook handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateEndpoint handles POST /v1/webhooks/endpoints
func (h *Handler) CreateEndpoint(c *gin.Context) {
	var input struct {
		URL    string   `json:"url" binding:"required,url"`
		Events []string `json:"events" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	tenantID, _ := middleware.GetTenantID(c)

	// Convert input events to Postgres Array string?
	// Our model uses `pg.StringArray` or similar? 
	// Let's check models.go: `Events string`. Again, likely implies JSON/Text array.
	// We need to fetch the tenant signing secret or verify we generate one.
	// Actually, `service.CreateEndpoint` takes a `models.WebhookEndpoint`.
	// We should generate a secret for the endpoint.
	
	secret, _ := h.service.GenerateEndpointSecret() // We need this method in service or just util

	endpoint := models.WebhookEndpoint{
		TenantID: tenantID,
		URL:      input.URL,
		// Events:   input.Events, // mismatch types
		Events:   toPostgresArray(input.Events), // Helper needed
		Secret:   secret,
		Enabled:  true,
	}

	created, err := h.service.CreateEndpoint(c.Request.Context(), endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "create_endpoint_failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// ListEndpoints handles GET /v1/webhooks/endpoints
func (h *Handler) ListEndpoints(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)
	
	// We need a List method in service
	endpoints, err := h.service.ListEndpoints(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "list_failed",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, endpoints)
}

// TODO: Move to shared util
func toPostgresArray(arr []string) string {
	if len(arr) == 0 {
		return "{}"
	}
	res := "{"
	for i, s := range arr {
		if i > 0 {
			res += ","
		}
		res += "\"" + s + "\""
	}
	res += "}"
	return res
}
