package usecases

import (
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/internal/usecase/auth/dtos"
	"github.com/paulozy/costurai/pkg"
)

type AuthDressmakerUseCase struct {
	DressMakerRepository database.DressmakerRepositoryInterface
}

type NewAuthDressmakerUseCaseInput struct {
	DressmakerRepository database.DressmakerRepositoryInterface
}

func NewDressmakerAuthenticationUseCase(repositories NewAuthDressmakerUseCaseInput) *AuthDressmakerUseCase {
	return &AuthDressmakerUseCase{
		DressMakerRepository: repositories.DressmakerRepository,
	}
}

func (useCase *AuthDressmakerUseCase) Execute(data dtos.AuthenticationInput) (dtos.AuthDressmakerOutput, pkg.Error) {
	dressmakerExists, err := useCase.DressMakerRepository.Exists(data.Email)
	if err != nil {
		return dtos.AuthDressmakerOutput{}, pkg.NewInternalServerError(err)
	}

	if !dressmakerExists {
		return dtos.AuthDressmakerOutput{}, pkg.NewInvalidCredentialsError()
	}

	dressmaker, err := useCase.DressMakerRepository.FindByEmail(data.Email)
	if err != nil {
		return dtos.AuthDressmakerOutput{}, pkg.NewInternalServerError(err)
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
