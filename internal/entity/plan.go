package entity

import "fmt"

type PlanType string
type PeriodicityType string

// Definitions of Plan

const (
	PlanTypeStandard PlanType = "standard"
	PlanTypePro      PlanType = "pro"
)

type Type struct {
	PlanType PlanType
}

const (
	MonthlyPeriodicity PeriodicityType = "monthly"
	YearlyPeriodicity  PeriodicityType = "yearly"
)

type Periodicity struct {
	PeriodicityType PeriodicityType
}

type Price struct {
	Amount    int32  `json:"amount"`
	Precision int8   `json:"precision"`
	Currency  string `json:"currency"`
}

type Plan struct {
	Name        Type        `json:"name"`
	Price       Price       `json:"price"`
	Periodicity Periodicity `json:"periodicity"`
}

// Definition of Builder

type PlanBuilder struct {
	plan Plan
	err  error
}

func NewPlanBuilder() *PlanBuilder {
	return &PlanBuilder{
		plan: Plan{},
	}
}

func (builder *PlanBuilder) WithType(planType PlanType) *PlanBuilder {
	builder.plan.Name = Type{PlanType: planType}
	return builder
}

func (builder *PlanBuilder) WithPeriodicity(p PeriodicityType) *PlanBuilder {
	builder.plan.Periodicity = Periodicity{PeriodicityType: p}
	return builder
}

func (builder *PlanBuilder) WithPrice(amount int32, precision int8, currency string) *PlanBuilder {
	builder.plan.Price = Price{
		Amount:    amount,
		Precision: precision,
		Currency:  currency,
	}
	return builder
}

func (b *PlanBuilder) ProStandard() *PlanBuilder {
	return b.WithType(PlanTypePro).
		WithPeriodicity(MonthlyPeriodicity).
		WithPrice(349, 2, "BRL")
}

func (b *PlanBuilder) ProPlan() *PlanBuilder {
	return b.WithType(PlanTypePro).
		WithPeriodicity(MonthlyPeriodicity).
		WithPrice(999, 2, "BRL")
}

func (b *PlanBuilder) Build() (Plan, error) {
	if b.err != nil {
		return Plan{}, b.err
	}
	if b.plan.Name.PlanType == "" {
		return Plan{}, fmt.Errorf("plan type is required")
	}
	if b.plan.Periodicity.PeriodicityType == "" {
		return Plan{}, fmt.Errorf("periodicity is required")
	}
	if b.plan.Price.Currency == "" {
		return Plan{}, fmt.Errorf("currency is required")
	}
	return b.plan, nil
}
