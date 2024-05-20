package usecases

import (
	"github.com/paulozy/costurai/internal/entity"
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/pkg"
)

type GetDressmakersByProximityInput struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Distance  int     `json:"distance"`
}

type GetDressmakersByProximityUseCase struct {
	DressMakerRepository database.DressmakerRepositoryInterface
}

func NewGetDressmakersByProximityUseCase(repo database.DressmakerRepositoryInterface) *GetDressmakersByProximityUseCase {
	return &GetDressmakersByProximityUseCase{
		DressMakerRepository: repo,
	}
}

func (useCase *GetDressmakersByProximityUseCase) Execute(data GetDressmakersByProximityInput) ([]entity.Dressmaker, pkg.Error) {
	dressmakers, err := useCase.DressMakerRepository.FindByProximity(data.Latitude, data.Longitude, data.Distance)

	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	return dressmakers, pkg.Error{}
}
