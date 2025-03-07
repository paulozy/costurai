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
	Dressmaker *entity.Dressmaker `json:"dressmaker,omitempty"`
	User       *entity.User       `json:"user,omitempty"`
	Token      string             `json:"token"`
}

func (useCase *AuthenticationUseCase) DressmakerExecute(data AuthenticationInput) (*AuthenticationOutput, pkg.Error) {
	dressmakerExists, err := useCase.DressMakerRepository.Exists(data.Email)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	if !dressmakerExists {
		return nil, pkg.NewInvalidCredentialsError()
	}

	dressmaker, err := useCase.DressMakerRepository.FindByEmail(data.Email)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	isValidPass := pkg.CompareHashAndPassword(dressmaker.Password, data.Password)
	if !isValidPass {
		return nil, pkg.NewInvalidCredentialsError()
	}

	token, err := pkg.GenerateToken(pkg.GenerateTokenInput{
		Issuer:  dressmaker.Name,
		Subject: dressmaker.ID,
	})
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	output := &AuthenticationOutput{
		Dressmaker: dressmaker,
		Token:      token,
	}

	return output, pkg.Error{}
}

func (useCase *AuthenticationUseCase) UserExecute(data AuthenticationInput) (*AuthenticationOutput, pkg.Error) {
	userExists, err := useCase.UserRepository.Exists(data.Email)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	if !userExists {
		return nil, pkg.NewInvalidCredentialsError()
	}

	user, err := useCase.UserRepository.FindByEmail(data.Email)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	isValidPass := pkg.CompareHashAndPassword(user.Password, data.Password)
	if !isValidPass {
		return nil, pkg.NewInvalidCredentialsError()
	}

	token, err := pkg.GenerateToken(pkg.GenerateTokenInput{
		Issuer:  user.Name,
		Subject: user.ID,
	})
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	output := &AuthenticationOutput{
		User:  user,
		Token: token,
	}

	return output, pkg.Error{}
}
