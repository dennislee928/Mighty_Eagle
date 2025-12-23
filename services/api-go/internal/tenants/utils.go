package tenants

import (
	"fmt"

	"github.com/google/uuid"
)

// parseUUID safely parses a UUID string
func parseUUID(s string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID: %w", err)
	}
	return id, nil
}

// ValidateTier checks if a tier is valid
func ValidateTier(tier string) bool {
	validTiers := map[string]bool{
		"lite":       true,
		"pro":        true,
		"enterprise": true,
	}
	return validTiers[tier]
}
