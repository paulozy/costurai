package usecases

import (
	"github.com/paulozy/costurai/internal/entity"
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/pkg"
)

type ShowDressMakerUseCase struct {
	DressMakerRepository database.DressmakerRepositoryInterface
}

func NewShowDressMakerUseCase(repo database.DressmakerRepositoryInterface) *ShowDressMakerUseCase {
	return &ShowDressMakerUseCase{
		DressMakerRepository: repo,
	}
}

type ShowDressMakerInput struct {
	ID string `json:"id"`
}

func (uc *ShowDressMakerUseCase) Execute(input ShowDressMakerInput) (*entity.Dressmaker, pkg.Error) {
	dressMaker, err := uc.DressMakerRepository.FindByID(input.ID)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	return dressMaker, pkg.Error{}
}
