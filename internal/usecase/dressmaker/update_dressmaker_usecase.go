package usecases

import (
	"github.com/paulozy/costurai/internal/entity"
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/pkg"
)

type UpdateDressMakerUseCase struct {
	DressmakerRepository database.DressmakerRepositoryInterface
}

func NewUpdateDressMakerUseCase(repo database.DressmakerRepositoryInterface) *UpdateDressMakerUseCase {
	return &UpdateDressMakerUseCase{
		DressmakerRepository: repo,
	}
}

func (uc *UpdateDressMakerUseCase) Execute(input entity.UpdateDressmakerInput) (*entity.Dressmaker, pkg.Error) {
	dressmaker, err := uc.DressmakerRepository.FindByID(input.ID)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	dressmaker.Update(input)

	err = uc.DressmakerRepository.Update(dressmaker)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	return dressmaker, pkg.Error{}
}
