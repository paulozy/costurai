package usecases

import (
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/pkg"
)

type AuthenticationUseCase struct {
	DressMakerRepository database.DressmakerRepositoryInterface
	UserRepository       database.UserRepositoryInterface
}

type NewAuthenticationUseCaseInput struct {
	DressMakerRepository database.DressmakerRepositoryInterface
	UserRepository       database.UserRepositoryInterface
}

func NewDressMakerAuthenticationUseCase(repositories NewAuthenticationUseCaseInput) *AuthenticationUseCase {
	return &AuthenticationUseCase{
		DressMakerRepository: repositories.DressMakerRepository,
		UserRepository:       repositories.UserRepository,
	}
}

type AuthenticationInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (useCase *AuthenticationUseCase) DressmakerExecute(data AuthenticationInput) (string, pkg.Error) {
	dressmakerExists, err := useCase.DressMakerRepository.Exists(data.Email)
	if err != nil {
		return "", pkg.NewInternalServerError(err)
	}

	if !dressmakerExists {
		return "", pkg.NewInvalidCredentialsError()
	}

	dressmaker, err := useCase.DressMakerRepository.FindByEmail(data.Email)
	if err != nil {
		return "", pkg.NewInternalServerError(err)
	}

	isValidPass := pkg.CompareHashAndPassword(dressmaker.Password, data.Password)
	if !isValidPass {
		return "", pkg.NewInvalidCredentialsError()
	}

	token, err := pkg.GenerateToken(pkg.GenerateTokenInput{
		Issuer:  dressmaker.Name,
		Subject: dressmaker.ID,
	})
	if err != nil {
		return "", pkg.NewInternalServerError(err)
	}

	return token, pkg.Error{}
}

func (useCase *AuthenticationUseCase) UserExecute(data AuthenticationInput) (string, pkg.Error) {
	userExists, err := useCase.UserRepository.Exists(data.Email)
	if err != nil {
		return "", pkg.NewInternalServerError(err)
	}

	if !userExists {
		return "", pkg.NewInvalidCredentialsError()
	}

	user, err := useCase.UserRepository.FindByEmail(data.Email)
	if err != nil {
		return "", pkg.NewInternalServerError(err)
	}

	isValidPass := pkg.CompareHashAndPassword(user.Password, data.Password)
	if !isValidPass {
		return "", pkg.NewInvalidCredentialsError()
	}

	token, err := pkg.GenerateToken(pkg.GenerateTokenInput{
		Issuer:  user.Name,
		Subject: user.ID,
	})
	if err != nil {
		return "", pkg.NewInternalServerError(err)
	}

	return token, pkg.Error{}
}