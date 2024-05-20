package usecases

import (
	"github.com/paulozy/costurai/internal/entity"
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/pkg"
)

type CreateDressMakerInput struct {
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	Name      string   `json:"name"`
	Contact   string   `json:"contact"`
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Services  []string `json:"services"`
}

type CreateDressMakerUseCase struct {
	DressMakerRepository database.DressmakerRepositoryInterface
}

func NewCreateDressMakerUseCase(repo database.DressmakerRepositoryInterface) *CreateDressMakerUseCase {
	return &CreateDressMakerUseCase{
		DressMakerRepository: repo,
	}
}

func (useCase *CreateDressMakerUseCase) Execute(data CreateDressMakerInput) (*entity.Dressmaker, pkg.Error) {
	validationError := validateInput(data)
	if validationError.Message != "" {
		return nil, validationError
	}

	dressmakerAlradyExists, _ := useCase.DressMakerRepository.FindByEmail(data.Email)

	if dressmakerAlradyExists != nil {
		return nil, pkg.NewEntityAlreadyExistsError("dressmaker")
	}

	location := entity.Location{
		Latitude:  data.Latitude,
		Longitude: data.Longitude,
	}

	dm, err := entity.NewDressmaker(data.Email, data.Password, data.Name, data.Contact, location, data.Services)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	err = useCase.DressMakerRepository.Create(dm)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	return dm, pkg.Error{}
}

func validateInput(data CreateDressMakerInput) pkg.Error {
	if data.Name == "" {
		return pkg.NewMissingFieldError("name")
	}

	if data.Email == "" {
		return pkg.NewMissingFieldError("email")
	}

	if data.Password == "" {
		return pkg.NewMissingFieldError("password")
	}

	if data.Contact == "" {
		return pkg.NewMissingFieldError("contact")
	}

	if data.Latitude == 0 {
		return pkg.NewMissingFieldError("latitude")
	}

	if data.Latitude < -85 || data.Latitude > 85 {
		return pkg.Error{
			Message: "latitude must be between -85 and 85",
			Status:  400,
		}
	}

	if data.Longitude == 0 {
		return pkg.NewMissingFieldError("longitude")
	}

	if data.Longitude < -180 || data.Longitude > 180 {
		return pkg.Error{
			Message: "longitude must be between -180 and 180",
			Status:  400,
		}
	}

	if len(data.Services) == 0 {
		return pkg.NewMissingFieldError("services")
	}

	return pkg.Error{}
}
