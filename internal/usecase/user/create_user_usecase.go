package usecases

import (
	"errors"

	"github.com/paulozy/costurai/internal/entity"
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/pkg"
)

type CreateUserUseCase struct {
	UserRepository database.UserRepositoryInterface
}

type CreateUserUseCaseInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func NewCreateUserUseCase(repo database.UserRepositoryInterface) *CreateUserUseCase {
	return &CreateUserUseCase{
		UserRepository: repo,
	}
}

func (useCase *CreateUserUseCase) Execute(data CreateUserUseCaseInput) (*entity.User, pkg.Error) {
	err := validateUserInput(data)
	if err != nil {
		return nil, pkg.NewBadRequestError(err)
	}

	userAlreadyExists, _ := useCase.UserRepository.FindByEmail(data.Email)
	if userAlreadyExists != nil {
		return nil, pkg.NewEntityAlreadyExistsError("user")
	}

	location := entity.Location{
		Latitude:  data.Latitude,
		Longitude: data.Longitude,
	}

	user, err := entity.NewUser(data.Email, data.Password, data.Name, location)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	err = useCase.UserRepository.Create(user)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	return user, pkg.Error{}
}

func validateUserInput(data CreateUserUseCaseInput) error {
	if data.Email == "" {
		return errors.New("email is required")
	}

	if data.Password == "" {
		return errors.New("password is required")
	}

	if data.Name == "" {
		return errors.New("name is required")
	}

	if data.Latitude == 0 {
		return errors.New("latitude is required")
	}

	if data.Latitude < -85 || data.Latitude > 85 {
		return errors.New("invalid latitude. latitude must be between -85 and 85")
	}

	if data.Longitude == 0 {
		return errors.New("longitude is required")
	}

	if data.Longitude < -180 || data.Longitude > 180 {
		return errors.New("invalid longitude. longitude must be between -180 and 180")
	}

	return nil
}