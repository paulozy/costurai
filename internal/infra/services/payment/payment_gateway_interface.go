package services

import "github.com/paulozy/costurai/internal/entity"

type PaymentPayload struct {
	Subscription *entity.Subscription
	Dressmaker   *entity.Dressmaker
	SuccessURL   string
	CancelURL    string
}

type PaymentGatewayServiceInterface interface {
	Pay(params PaymentPayload) (string, error)
}
