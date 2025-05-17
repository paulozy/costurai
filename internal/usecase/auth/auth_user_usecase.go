package usecases

import (
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/internal/usecase/auth/dtos"
	"github.com/paulozy/costurai/pkg"
)

type AuthUserUseCase struct {
	UserRepository database.UserRepositoryInterface
}

type NewAuthUserUseCaseInput struct {
	UserRepository database.UserRepositoryInterface
}

func NewUserAuthUseCase(repositories NewAuthUserUseCaseInput) *AuthUserUseCase {
	return &AuthUserUseCase{
		UserRepository: repositories.UserRepository,
	}
}

func (useCase *AuthUserUseCase) UserExecute(data dtos.AuthenticationInput) (dtos.AuthUserOutput, pkg.Error) {
	userExists, err := useCase.UserRepository.Exists(data.Email)
	if err != nil {
		return dtos.AuthUserOutput{}, pkg.NewInternalServerError(err)
	}

	if !userExists {
		return dtos.AuthUserOutput{}, pkg.NewInvalidCredentialsError()
	}

	user, err := useCase.UserRepository.FindByEmail(data.Email)
	if err != nil {
		return dtos.AuthUserOutput{}, pkg.NewInternalServerError(err)
	}

	isValidPass := pkg.CompareHashAndPassword(user.Password, data.Password)
	if !isValidPass {
		return dtos.AuthUserOutput{}, pkg.NewInvalidCredentialsError()
	}

	token, err := pkg.GenerateToken(pkg.GenerateTokenInput{
		Issuer:  user.Name,
		Subject: user.ID,
	})
	if err != nil {
		return dtos.AuthUserOutput{}, pkg.NewInternalServerError(err)
	}

	response := dtos.AuthUserOutput{
		Token: token,
		User:  *user,
	}

	return response, pkg.Error{}
}
