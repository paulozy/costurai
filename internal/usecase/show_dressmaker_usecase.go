package usecases

import (
	"github.com/paulozy/costurai/internal/entity"
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/pkg"
)

type ShowDressmakerUseCase struct {
	DressMakerRepository database.DressmakerRepositoryInterface
}

func NewShowDressmakerUseCase(repo database.DressmakerRepositoryInterface) *ShowDressmakerUseCase {
	return &ShowDressmakerUseCase{
		DressMakerRepository: repo,
	}
}

type ShowDressmakerInput struct {
	ID string `json:"id"`
}

func (useCase *ShowDressmakerUseCase) Execute(input ShowDressmakerInput) (*entity.Dressmaker, pkg.Error) {
	dressmaker, err := useCase.DressMakerRepository.FindByID(input.ID)
	if err != nil {
		return nil, pkg.NewNotFoundError("dressmaker")
	}

	return dressmaker, pkg.Error{}
}
