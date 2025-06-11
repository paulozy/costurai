package usecases

import (
	"github.com/paulozy/costurai/configs"
	"github.com/paulozy/costurai/internal/entity"
	"github.com/paulozy/costurai/internal/infra/database"
	services "github.com/paulozy/costurai/internal/infra/services/payment"
	"github.com/paulozy/costurai/pkg"
)

type CreateSubscriptionUseCase struct {
	SubscriptionRepository database.SubscriptionRepositoryInterface
	DressmakerRepository   database.DressmakerRepositoryInterface
	PaymentGatewayService  services.PaymentGatewayServiceInterface
	Configs                *configs.Config
}

func NewCreateSubscriptionUseCase(
	subRepo database.SubscriptionRepositoryInterface,
	dmRepo database.DressmakerRepositoryInterface,
	paymentGatewayService services.PaymentGatewayServiceInterface,
	cfg *configs.Config,
) *CreateSubscriptionUseCase {
	return &CreateSubscriptionUseCase{
		SubscriptionRepository: subRepo,
		DressmakerRepository:   dmRepo,
		PaymentGatewayService:  paymentGatewayService,
		Configs:                cfg,
	}
}

type CreateSubscriptionInput struct {
	DressmakerID    string                 `json:"dressmakerId"`
	PlanType        entity.PlanType        `json:"planType"`
	PeriodicityType entity.PeriodicityType `json:"periodicityType"`
}

func (uc *CreateSubscriptionUseCase) Execute(input CreateSubscriptionInput) (*entity.Subscription, pkg.Error) {
	dressmaker, err := uc.DressmakerRepository.FindByID(input.DressmakerID)
	if err != nil {
		return nil, pkg.Error{
			Error:   err.Error(),
			Message: "Error on getting dressmaker",
		}
	}

	plan, err := uc.getPlan(input.PlanType, input.PeriodicityType)
	if err != nil {
		return nil, pkg.Error{
			Error:   err.Error(),
			Message: "Error creating a new Plan",
		}
	}

	subscription, err := entity.NewSubscription(dressmaker.ID, *plan)
	if err != nil {
		return nil, pkg.Error{
			Error:   err.Error(),
			Message: "Error creating a new Subscription",
		}
	}

	paymentPayload := services.PaymentPayload{
		Subscription: subscription,
		Dressmaker:   dressmaker,
		SuccessURL:   uc.Configs.PaymentSuccessRedirectURL,
		CancelURL:    uc.Configs.PaymentCancelRedirectURL,
	}

	redirectURL, err := uc.PaymentGatewayService.Pay(paymentPayload)
	if err != nil {
		return nil, pkg.Error{
			Error:   err.Error(),
			Message: "Error on creating payment intent",
		}
	}

	subscription.PaymentURL = &redirectURL
	err = uc.SubscriptionRepository.Create(subscription)
	if err != nil {
		return nil, pkg.Error{
			Error:   err.Error(),
			Message: "Error saving subscription",
		}
	}

	dressmaker.SubscriptionId = &subscription.ID
	err = uc.DressmakerRepository.Update(dressmaker)
	if err != nil {
		return nil, pkg.Error{
			Error:   err.Error(),
			Message: "Error saving dressmaker",
		}
	}

	return subscription, pkg.Error{}
}

func (uc *CreateSubscriptionUseCase) getPlan(planType entity.PlanType, periodicity entity.PeriodicityType) (*entity.Plan, error) {
	plannerPrice := pkg.NewPlannerPrice()
	planBuilder := entity.NewPlanBuilder()
	prePlan := planBuilder.WithType(planType).WithPeriodicity(periodicity)

	isStandardMonthly := planType == entity.PlanTypeStandard &&
		periodicity == entity.MonthlyPeriodicity

	isStandardYearly := planType == entity.PlanTypeStandard &&
		periodicity == entity.YearlyPeriodicity

	isProMonthly := planType == entity.PlanTypePro &&
		periodicity == entity.MonthlyPeriodicity

	if isStandardMonthly {
		prePlan.WithPrice(plannerPrice.MonthlyStandard)
	} else if isStandardYearly {
		prePlan.WithPrice(plannerPrice.YearlyStandard)
	} else if isProMonthly {
		prePlan.WithPrice(plannerPrice.MonthlyPro)
	} else {
		prePlan.WithPrice(plannerPrice.YearlyPro)
	}

	plan, err := prePlan.Build()
	if err != nil {
		return nil, err
	}

	return &plan, nil
}
