package services

import (
	"fmt"

	"github.com/paulozy/costurai/internal/entity"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/price"
	"github.com/stripe/stripe-go/v82/product"
)

type StripeService struct {
}

func NewStripeService() *StripeService {
	return &StripeService{}
}

func InitStripe(secretKey string) {
	stripe.Key = secretKey
}

var WebhookSecret string

func InitWebhook(secret string) {
	WebhookSecret = secret
}

func (s *StripeService) Pay(params PaymentPayload) (string, error) {
	priceID, err := s.getStripePriceID(params.Subscription.Plan)
	if err != nil {
		return s.createInlineCheckoutSession(params)
	}

	return s.createCheckoutSessionWithPriceID(params, priceID)
}

func (s *StripeService) createCheckoutSessionWithPriceID(req PaymentPayload, priceID string) (string, error) {
	params := &stripe.CheckoutSessionParams{
		CustomerEmail: stripe.String(req.Dressmaker.Email),
		SuccessURL:    stripe.String(req.SuccessURL),
		CancelURL:     stripe.String(req.CancelURL),
		Mode:          stripe.String("subscription"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Metadata: map[string]string{
			"subscription_id": req.Subscription.ID,
			"dressmaker_id":   req.Dressmaker.ID,
		},
	}

	session, err := session.New(params)
	if err != nil {
		return "", err
	}

	return session.URL, nil
}

func (s *StripeService) createInlineCheckoutSession(req PaymentPayload) (string, error) {
	plan := req.Subscription.Plan
	productName := s.getProductName(plan)
	interval := s.getInterval(plan.Periodicity)

	params := &stripe.CheckoutSessionParams{
		CustomerEmail: stripe.String(req.Dressmaker.Email),
		SuccessURL:    stripe.String(req.SuccessURL),
		CancelURL:     stripe.String(req.CancelURL),
		Mode:          stripe.String("subscription"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:   stripe.String(plan.Price.Currency),
					UnitAmount: stripe.Int64(int64(plan.Price.Amount)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(productName),
					},
					Recurring: &stripe.CheckoutSessionLineItemPriceDataRecurringParams{
						Interval: stripe.String(interval),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
		Metadata: map[string]string{
			"subscription_id": req.Subscription.ID,
			"dressmaker_id":   req.Dressmaker.ID,
		},
	}

	session, err := session.New(params)
	if err != nil {
		return "", err
	}

	return session.URL, nil
}

func (s *StripeService) getStripePriceID(plan entity.Plan) (string, error) {
	productName := s.getProductName(plan)
	interval := string(plan.Periodicity) // "monthly" or "yearly"

	productList := &stripe.ProductListParams{
		Active: stripe.Bool(true),
	}
	i := product.List(productList)

	for i.Next() {
		p := i.Product()
		if p.Name == productName {
			priceList := &stripe.PriceListParams{
				Product: stripe.String(p.ID),
				Recurring: &stripe.PriceListRecurringParams{
					Interval: stripe.String(interval),
				},
			}
			pi := price.List(priceList)
			for pi.Next() {
				return pi.Price().ID, nil
			}
		}
	}

	return "", fmt.Errorf("price ID not found")
}

func (s *StripeService) getProductName(plan entity.Plan) string {
	var productName string

	isStandardMonthly := plan.Name == entity.PlanTypeStandard &&
		plan.Periodicity == entity.MonthlyPeriodicity

	isStandardYearly := plan.Name == entity.PlanTypeStandard &&
		plan.Periodicity == entity.YearlyPeriodicity

	isProMonthly := plan.Name == entity.PlanTypePro &&
		plan.Periodicity == entity.MonthlyPeriodicity

	if isStandardMonthly {
		productName = "Inscrição Standard - Mensal"
	} else if isStandardYearly {
		productName = "Inscrição Standard - Anual"
	} else if isProMonthly {
		productName = "Inscrição Pro - Mensal"
	} else {
		productName = "Inscrição Pro - Anual"
	}

	return productName
}

func (s *StripeService) getInterval(periodicity entity.PeriodicityType) string {
	switch periodicity {
	case entity.MonthlyPeriodicity:
		return "month"
	case entity.YearlyPeriodicity:
		return "year"
	default:
		return "month"
	}
}
