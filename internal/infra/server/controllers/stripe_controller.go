package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"

	"github.com/paulozy/costurai/internal/entity"
	"github.com/paulozy/costurai/internal/infra/database"
	paymentServices "github.com/paulozy/costurai/internal/infra/services/payment"
)

type StripeController struct {
	subscriptionRepository database.SubscriptionRepositoryInterface
	dressmakerRepository   database.DressmakerRepositoryInterface
}

func NewStripeController(
	subRepo database.SubscriptionRepositoryInterface,
	dmRepo database.DressmakerRepositoryInterface,
) *StripeController {
	return &StripeController{
		subscriptionRepository: subRepo,
		dressmakerRepository:   dmRepo,
	}
}

func (sc *StripeController) HandleWebhook(c *gin.Context) {
	payload, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("could not read request body: %v", err))
		return
	}
	sigHeader := c.GetHeader("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, sigHeader, paymentServices.WebhookSecret)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("webhook signature verification failed: %v", err))
		return
	}
	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		var sess stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &sess); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("could not parse webhook JSON: %v", err))
			return
		}
		subID := sess.Metadata["subscription_id"]
		sub, err := sc.subscriptionRepository.FindByID(subID)
		if err != nil {
			c.String(http.StatusNotFound, fmt.Sprintf("subscription not found: %v", err))
			return
		}
		sub.Status = entity.StatusActive
		if sess.Subscription != nil {
			gatewayID := sess.Subscription.ID
			sub.GatewayId = &gatewayID
		}
		if err := sc.subscriptionRepository.Update(sub); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("could not update subscription: %v", err))
			return
		}
	default:
	}
	c.String(http.StatusOK, "success")
}
