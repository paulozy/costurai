package usecases

import (
	"github.com/paulozy/costurai/internal/entity"
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

type AuthenticationOutput struct {
	Token string            `json:"token"`
	User  entity.Dressmaker `json:"user"`
}

func (useCase *AuthenticationUseCase) DressmakerExecute(data AuthenticationInput) (AuthenticationOutput, pkg.Error) {
	dressmakerExists, err := useCase.DressMakerRepository.Exists(data.Email)
	if err != nil {
		return AuthenticationOutput{}, pkg.NewInternalServerError(err)
	}

	if !dressmakerExists {
		return AuthenticationOutput{}, pkg.NewInvalidCredentialsError()
	}

	dressmaker, err := useCase.DressMakerRepository.FindByEmail(data.Email)
	if err != nil {
		return AuthenticationOutput{}, pkg.NewInternalServerError(err)
	}

	isValidPass := pkg.CompareHashAndPassword(dressmaker.Password, data.Password)
	if !isValidPass {
		return AuthenticationOutput{}, pkg.NewInvalidCredentialsError()
	}

	token, err := pkg.GenerateToken(pkg.GenerateTokenInput{
		Issuer:  dressmaker.Name,
		Subject: dressmaker.ID,
	})
	if err != nil {
		return AuthenticationOutput{}, pkg.NewInternalServerError(err)
	}

	response := AuthenticationOutput{
		Token: token,
		User:  *dressmaker,
	}

	return response, pkg.Error{}
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
