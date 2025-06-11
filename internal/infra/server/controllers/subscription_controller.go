package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/paulozy/costurai/internal/infra/database"
	usecases "github.com/paulozy/costurai/internal/usecase/subscription"
)

type SubscriptionController struct {
	dressmakerRepository      database.DressmakerRepositoryInterface
	subscriptionRepository    database.SubscriptionRepositoryInterface
	createSubscriptionUseCase *usecases.CreateSubscriptionUseCase
}

type SubscriptionUseCasesInput struct {
	CreateSubscriptionUseCase *usecases.CreateSubscriptionUseCase
}

func NewSubscriptionController(
	dmRepo database.DressmakerRepositoryInterface,
	subRepo database.SubscriptionRepositoryInterface,
	usecases SubscriptionUseCasesInput,
) *SubscriptionController {
	return &SubscriptionController{
		dressmakerRepository:      dmRepo,
		subscriptionRepository:    subRepo,
		createSubscriptionUseCase: usecases.CreateSubscriptionUseCase,
	}
}

func (sc *SubscriptionController) CreateSubscription(c *gin.Context) {
	var input usecases.CreateSubscriptionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	checkoutURL, err := sc.createSubscriptionUseCase.Execute(input)
	if err.Message != "" {
		c.JSON(err.Status, gin.H{"error": err.Message, "reason": err.Error})
		return
	}

	c.JSON(201, gin.H{"data": checkoutURL})
}
