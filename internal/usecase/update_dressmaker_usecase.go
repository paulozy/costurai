package usecases

import (
	"github.com/paulozy/costurai/internal/entity"
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/pkg"
)

type UpdateDressMakerUseCase struct {
	DressMakerRepository database.DressmakerRepositoryInterface
}

func NewUpdateDressMakerUseCase(repo database.DressmakerRepositoryInterface) *UpdateDressMakerUseCase {
	return &UpdateDressMakerUseCase{
		DressMakerRepository: repo,
	}
}

type UpdateDressMakerInput struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Contact   string  `json:"contact"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (uc *UpdateDressMakerUseCase) Execute(input UpdateDressMakerInput) (*entity.Dressmaker, pkg.Error) {
	dressMaker, err := uc.DressMakerRepository.FindByID(input.ID)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	dressMaker.Name = input.Name
	dressMaker.Contact = input.Contact
	dressMaker.Location.Latitude = input.Latitude
	dressMaker.Location.Longitude = input.Longitude

	err = uc.DressMakerRepository.Update(dressMaker)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	return dressMaker, pkg.Error{}
}
