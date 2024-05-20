package usecases

import (
	"fmt"

	"github.com/paulozy/costurai/internal/entity"
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/pkg"
)


type GetDressmakersByServicesInput struct {
	Services string `json:"services"`
}

type GetDressmakersByServicesUseCase struct {
	DressMakerRepository database.DressmakerRepositoryInterface
}

func NewGetDressmakersByServicesUseCase(repo database.DressmakerRepositoryInterface) *GetDressmakersByServicesUseCase {
	return &GetDressmakersByServicesUseCase{
		DressMakerRepository: repo,
	}
}

func (useCase *GetDressmakersByServicesUseCase) Execute(data GetDressmakersByServicesInput) ([]entity.Dressmaker, pkg.Error) {
	println("GetDressmakersByServicesUseCase.Execute", data.Services)

	dressmakers, err := useCase.DressMakerRepository.FindByServices(data.Services)

	if err != nil {
		fmt.Println("Error on GetDressmakersByServicesUseCase.Execute", err)
		return nil, pkg.NewInternalServerError(err)
	}

	return dressmakers, pkg.Error{}
}