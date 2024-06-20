package usecases

import (
	"fmt"

	"github.com/paulozy/costurai/internal/entity"
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/pkg"
)

type GetDressmakersUseCase struct {
	DressMakerRepository database.DressmakerRepositoryInterface
}

func NewGetDressmakersUseCase(repo database.DressmakerRepositoryInterface) *GetDressmakersUseCase {
	return &GetDressmakersUseCase{
		DressMakerRepository: repo,
	}
}

type GetDressmakersInput struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Distance  int     `json:"distance"`

	Services string `json:"services"`
}

func (useCase *GetDressmakersUseCase) Execute(input GetDressmakersInput) ([]entity.Dressmaker, pkg.Error) {
	var searchParams database.GetDressmakersParams

	switch {
	case input.Services != "":
		searchParams = database.GetDressmakersParams{
			Services: input.Services,
		}
	case input.Latitude != 0 && input.Longitude != 0 && input.Distance != 0:
		searchParams = database.GetDressmakersParams{
			Latitude:  input.Latitude,
			Longitude: input.Longitude,
			Distance:  input.Distance,
		}
	default:
		searchParams = database.GetDressmakersParams{
			Default: true,
		}
	}

	fmt.Println(searchParams)

	dressmakers, err := useCase.DressMakerRepository.Find(searchParams)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	return dressmakers, pkg.Error{}
}
