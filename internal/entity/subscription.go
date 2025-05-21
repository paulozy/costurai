package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID           string      `json:"id"`
	DressmakerID string      `json:"dressmakerId"`
	PlanType     PlanType    `json:"planType"` // padrão fixo: "standard", "pro"
	Price        Price       `json:"price"`
	Periodicity  Periodicity `json:"periodicity"`

	StartedAt  time.Time  `json:"startedAt"`
	ExpiresAt  *time.Time `json:"expiresAt,omitempty"`
	CanceledAt *time.Time `json:"canceledAt,omitempty"`
	Active     bool       `json:"active"`
	GraceUntil *time.Time `json:"graceUntil,omitempty"` // até quando mantém acesso
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

func NewSubscription(dressmakerID string, plan Plan) (*Subscription, error) {
	if plan.Name.PlanType == "" {
		return nil, fmt.Errorf("plan type is required")
	}

	duration, err := durationForPeriodicity(plan.Periodicity.PeriodicityType)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	expires := now.Add(duration)

	return &Subscription{
		ID:           uuid.New().String(),
		DressmakerID: dressmakerID,
		PlanType:     plan.Name.PlanType,
		Price:        plan.Price,
		Periodicity:  plan.Periodicity,
		StartedAt:    now,
		ExpiresAt:    &expires,
		Active:       true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (s *Subscription) IsActive() bool {
	if !s.Active {
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
	s.Active = false
	s.CanceledAt = &now

	if gracePeriodDays > 0 {
		grace := now.AddDate(0, 0, gracePeriodDays)
		s.GraceUntil = &grace
	}
}

func (s *Subscription) Renew() error {
	if s.Periodicity.PeriodicityType == "" {
		return fmt.Errorf("periodicity type undefined")
	}

	var duration time.Duration
	switch s.Periodicity.PeriodicityType {
	case MonthlyPeriodicity:
		duration = 30 * 24 * time.Hour
	case YearlyPeriodicity:
		duration = 365 * 24 * time.Hour
	default:
		return fmt.Errorf("unsupported periodicity")
	}

	now := time.Now()
	s.StartedAt = now
	expires := now.Add(duration)
	s.ExpiresAt = &expires
	s.Active = true
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
