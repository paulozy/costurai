package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusPending  Status = "pending"
	StatusCanceled Status = "canceled"
)

type Subscription struct {
	ID           string `json:"id"`
	DressmakerID string `json:"dressmakerId"`
	Plan         Plan   `json:"plan"` // padrão fixo: "standard", "pro"
	Status       Status `json:"status"`

	StartedAt  *time.Time `json:"startedAt"`
	ExpiresAt  *time.Time `json:"expiresAt,omitempty"`
	CanceledAt *time.Time `json:"canceledAt,omitempty"`
	GraceUntil *time.Time `json:"graceUntil,omitempty"` // até quando mantém acesso
	GatewayId  *string    `json:"gatewayId,omitempty"`
	PaymentURL *string    `json:"paymentURL,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewSubscription(dressmakerID string, plan Plan) (*Subscription, error) {
	if plan.Name == "" {
		return nil, fmt.Errorf("plan type is required")
	}

	duration, err := durationForPeriodicity(plan.Periodicity)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	expires := now.Add(duration)

	return &Subscription{
		ID:           uuid.New().String(),
		DressmakerID: dressmakerID,
		Plan:         plan,
		StartedAt:    &now,
		ExpiresAt:    &expires,
		Status:       StatusPending,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (s *Subscription) IsActive() bool {
	if s.Status != StatusActive {
		return false
	}
	if s.ExpiresAt != nil && s.ExpiresAt.Before(time.Now()) {
		return false
	}
	return true
}

func (s *Subscription) IsInGracePeriod() bool {
	if s.GraceUntil == nil {
		return false
	}
	now := time.Now()
	return now.Before(*s.GraceUntil)
}

func (s *Subscription) HasExpired() bool {
	return s.ExpiresAt != nil && time.Now().After(*s.ExpiresAt)
}

func (s *Subscription) Cancel(gracePeriodDays int) {
	now := time.Now()
	s.Status = StatusCanceled
	s.CanceledAt = &now

	if gracePeriodDays > 0 {
		grace := now.AddDate(0, 0, gracePeriodDays)
		s.GraceUntil = &grace
	}
}

func (s *Subscription) Renew() error {
	if s.Plan.Periodicity == "" {
		return fmt.Errorf("periodicity type undefined")
	}

	var duration time.Duration
	switch s.Plan.Periodicity {
	case MonthlyPeriodicity:
		duration = 30 * 24 * time.Hour
	case YearlyPeriodicity:
		duration = 365 * 24 * time.Hour
	default:
		return fmt.Errorf("unsupported periodicity")
	}

	now := time.Now()
	s.StartedAt = &now
	expires := now.Add(duration)
	s.ExpiresAt = &expires
	s.Status = StatusActive
	s.CanceledAt = nil
	s.GraceUntil = nil
	return nil
}

func durationForPeriodicity(p PeriodicityType) (time.Duration, error) {
	switch p {
	case MonthlyPeriodicity:
		return 30 * 24 * time.Hour, nil
	case YearlyPeriodicity:
		return 365 * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unsupported periodicity: %s", p)
	}
}
