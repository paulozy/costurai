package usecases

import (
	"github.com/paulozy/costurai/internal/entity"
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/pkg"
)

type ShowDressmakerUseCase struct {
	DressmakerRepository   database.DressmakerRepositoryInterface
	SubscriptionRepository database.SubscriptionRepositoryInterface
}

func NewShowDressmakerUseCase(
	repo database.DressmakerRepositoryInterface,
	subscriptionRepo database.SubscriptionRepositoryInterface,
) *ShowDressmakerUseCase {
	return &ShowDressmakerUseCase{
		DressmakerRepository:   repo,
		SubscriptionRepository: subscriptionRepo,
	}
}

type ShowDressmakerInput struct {
	ID string `json:"id"`
}

func (uc *ShowDressmakerUseCase) Execute(input ShowDressmakerInput) (*entity.Dressmaker, pkg.Error) {
	dressmaker, err := uc.DressmakerRepository.FindByID(input.ID)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	if dressmaker.SubscriptionId != nil {
		subscription, err := uc.SubscriptionRepository.FindByID(*dressmaker.SubscriptionId)
		if err != nil {
			return nil, pkg.NewInternalServerError(err)
		}
		dressmaker.Subscription = subscription
	}

	return dressmaker, pkg.Error{}
}
