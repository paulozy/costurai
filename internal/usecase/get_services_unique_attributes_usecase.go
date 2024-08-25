package usecases

import (
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/pkg"
)

type GetServicesUniqueAttributesUseCase struct {
	DressMakerRepository database.DressmakerRepositoryInterface
}

func NewGetServicesUniqueAttributesUseCase(repo database.DressmakerRepositoryInterface) *GetServicesUniqueAttributesUseCase {
	return &GetServicesUniqueAttributesUseCase{
		DressMakerRepository: repo,
	}
}

func (usecase *GetServicesUniqueAttributesUseCase) Execute() ([]string, pkg.Error) {

	dressmakersServices, err := usecase.DressMakerRepository.GetServices()
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	return dressmakersServices, pkg.Error{}
}
