package usecases

import (
	"github.com/paulozy/costurai/internal/entity"
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/pkg"
)

type CreateDressMakerUseCase struct {
	DressmakerRepository database.DressmakerRepositoryInterface
}

func NewCreateDressMakerUseCase(repo database.DressmakerRepositoryInterface) *CreateDressMakerUseCase {
	return &CreateDressMakerUseCase{
		DressmakerRepository: repo,
	}
}

func (useCase *CreateDressMakerUseCase) Execute(data entity.CreateDressmakerInput) (*entity.Dressmaker, pkg.Error) {
	validationError := validateInput(data)
	if validationError.Message != "" {
		return nil, validationError
	}

	dressmakerAlradyExists, _ := useCase.DressmakerRepository.FindByEmail(data.Email)

	if dressmakerAlradyExists != nil {
		return nil, pkg.NewEntityAlreadyExistsError("dressmaker")
	}

	dressmaker, err := entity.NewDressmaker(data)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	err = useCase.DressmakerRepository.Create(dressmaker)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	return dressmaker, pkg.Error{}
}

func validateInput(data entity.CreateDressmakerInput) pkg.Error {
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

	if (data.Address == entity.Address{}) {
		return pkg.NewMissingFieldError("address")
	}

	if data.Address.Street == "" {
		return pkg.NewMissingFieldError("address.street")
	}

	if data.Address.Number == "" {
		return pkg.NewMissingFieldError("address.number")
	}

	if data.Address.Neighborhood == "" {
		return pkg.NewMissingFieldError("address.neighborhood")
	}

	if data.Address.City == "" {
		return pkg.NewMissingFieldError("address.city")
	}

	if data.Address.State == "" {
		return pkg.NewMissingFieldError("address.state")
	}

	if data.Address.Location.Latitude == 0 {
		return pkg.NewMissingFieldError("latitude")
	}

	if data.Address.Location.Latitude < -85 || data.Address.Location.Latitude > 85 {
		return pkg.Error{
			Message: "latitude must be between -85 and 85",
			Status:  400,
		}
	}

	if data.Address.Location.Longitude == 0 {
		return pkg.NewMissingFieldError("longitude")
	}

	if data.Address.Location.Longitude < -180 || data.Address.Location.Longitude > 180 {
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
