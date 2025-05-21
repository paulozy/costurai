package usecases

import (
	"github.com/paulozy/costurai/internal/entity"
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/pkg"
	"github.com/paulozy/costurai/pkg/paginator"
)

type GetDressmakersByProximityInput struct {
	Latitude  float64 `form:"latitude"`
	Longitude float64 `form:"longitude"`
	Distance  int     `form:"distance"`

	Limit int64 `form:"limit"`
	Page  int64 `form:"page"`
}

type GetDressmakersByProximityUseCase struct {
	DressMakerRepository database.DressmakerRepositoryInterface
}

func NewGetDressmakersByProximityUseCase(repo database.DressmakerRepositoryInterface) *GetDressmakersByProximityUseCase {
	return &GetDressmakersByProximityUseCase{
		DressMakerRepository: repo,
	}
}

type GetDressmakersByProximityOuput struct {
	*paginator.Paginate[entity.Dressmaker]
}

func (useCase *GetDressmakersByProximityUseCase) Execute(data GetDressmakersByProximityInput) (*GetDressmakersByProximityOuput, pkg.Error) {
	dressmakers, err := useCase.DressMakerRepository.FindByProximity(data.Latitude, data.Longitude, data.Distance)

	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	offset := paginator.GetOffset(data.Limit, data.Page, dressmakers)
	paginatedItems := dressmakers[offset.Start:offset.End]

	response := &GetDressmakersByProximityOuput{
		Paginate: &paginator.Paginate[entity.Dressmaker]{
			Items:          &paginatedItems,
			PaginationInfo: paginator.NewPaginatation(data.Limit, data.Page, int64(len(dressmakers))),
		},
	}

	return response, pkg.Error{}
}
