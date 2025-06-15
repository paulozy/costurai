package usecases

import (
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/internal/usecase/auth/dtos"
	"github.com/paulozy/costurai/pkg"
)

type AuthDressmakerUseCase struct {
	DressmakerRepository   database.DressmakerRepositoryInterface
	SubscriptionRepository database.SubscriptionRepositoryInterface
}

type NewAuthDressmakerUseCaseInput struct {
	DressmakerRepository   database.DressmakerRepositoryInterface
	SubscriptionRepository database.SubscriptionRepositoryInterface
}

func NewDressmakerAuthenticationUseCase(repositories NewAuthDressmakerUseCaseInput) *AuthDressmakerUseCase {
	return &AuthDressmakerUseCase{
		DressmakerRepository:   repositories.DressmakerRepository,
		SubscriptionRepository: repositories.SubscriptionRepository,
	}
}

func (uc *AuthDressmakerUseCase) Execute(data dtos.AuthenticationInput) (dtos.AuthDressmakerOutput, pkg.Error) {
	dressmakerExists, err := uc.DressmakerRepository.Exists(data.Email)
	if err != nil {
		return dtos.AuthDressmakerOutput{}, pkg.NewInternalServerError(err)
	}

	if !dressmakerExists {
		return dtos.AuthDressmakerOutput{}, pkg.NewInvalidCredentialsError()
	}

	dressmaker, err := uc.DressmakerRepository.FindByEmail(data.Email)
	if err != nil {
		return dtos.AuthDressmakerOutput{}, pkg.NewInternalServerError(err)
	}

	if dressmaker.SubscriptionId != nil {
		subscription, err := uc.SubscriptionRepository.FindByID(*dressmaker.SubscriptionId)
		if err != nil {
			return dtos.AuthDressmakerOutput{}, pkg.NewInternalServerError(err)
		}
		dressmaker.Subscription = subscription
	}

	isValidPass := pkg.CompareHashAndPassword(dressmaker.Password, data.Password)
	if !isValidPass {
		return dtos.AuthDressmakerOutput{}, pkg.NewInvalidCredentialsError()
	}

	token, err := pkg.GenerateToken(pkg.GenerateTokenInput{
		Issuer:  dressmaker.Name,
		Subject: dressmaker.ID,
	})
	if err != nil {
		return dtos.AuthDressmakerOutput{}, pkg.NewInternalServerError(err)
	}

	response := dtos.AuthDressmakerOutput{
		Token: token,
		User:  *dressmaker,
	}

	return response, pkg.Error{}
}
