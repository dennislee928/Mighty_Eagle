package reputation

import (
	"time"

	"github.com/dennislee928/mighty-eagle/api-go/internal/models"
)

// ScoreComponents breakdown of how score was calculated
type ScoreComponents struct {
	BaseScore        int `json:"base_score"`
	Verification     int `json:"verification_bonus"`
	AccountAge       int `json:"account_age_bonus"`
	History          int `json:"history_bonus"` // or penalty
	Signals          int `json:"signals_bonus"` // or penalty
}

// Scorer implements the deterministic scoring algorithm
type Scorer struct {
	// Configuration for weights
	VerifiedWeight    int
	AgeWeightPerMonth int
	MaxAgeBonus       int
}

// NewScorer creates a default scorer
func NewScorer() *Scorer {
	return &Scorer{
		VerifiedWeight:    40,
		AgeWeightPerMonth: 1, // 1 point per month
		MaxAgeBonus:       30, // Cap age bonus at 30 points (2.5 years)
	}
}

// CalculateScore computes the reputation score for a subject
func (s *Scorer) CalculateScore(verifications []models.PersonaVerification, history []interface{}) (float64, ScoreComponents) {
	components := ScoreComponents{
		BaseScore: 0,
	}

	// 1. Verification Status (+40)
	// Check if any verification is active and verified
	isVerified := false
	var earliestVerification *time.Time

	for _, v := range verifications {
		if v.Status == "verified" {
			// Check expiration
			if v.ExpiresAt == nil || v.ExpiresAt.After(time.Now()) {
				isVerified = true
				if earliestVerification == nil || v.CreatedAt.Before(*earliestVerification) {
					createdAt := v.CreatedAt
					earliestVerification = &createdAt
				}
			}
		}
	}

	if isVerified {
		components.Verification = s.VerifiedWeight
	}

	// 2. Account Age (0-30 verification age)
	// For API MVP, we might rely on the earliest verification date as "account age" in our system
	// unless we are passed an account created_at date from the tenant.
	// Let's assume age starts from first verification for now.
	if earliestVerification != nil {
		age := time.Since(*earliestVerification)
		months := int(age.Hours() / 24 / 30)
		ageBonus := months * s.AgeWeightPerMonth
		if ageBonus > s.MaxAgeBonus {
			ageBonus = s.MaxAgeBonus
		}
		components.AccountAge = ageBonus
	}

	// 3. Dispute Signals / History (-20 to +30)
	// TODO: Signal integration. For M3 MVP, we assume neutral (0).
	components.Signals = 0

	// Calculate Total
	total := components.BaseScore + components.Verification + components.AccountAge + components.History + components.Signals

	// Cap filter (0-100)
	if total < 0 {
		total = 0
	}
	if total > 100 {
		total = 100
	}

	return float64(total), components
}
