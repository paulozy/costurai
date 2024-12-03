package usecases

import (
	"fmt"

	"github.com/paulozy/costurai/internal/entity"
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/pkg"
	"github.com/paulozy/costurai/pkg/paginator"
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
	Latitude  float64 `form:"latitude"`
	Longitude float64 `form:"longitude"`
	Distance  int     `form:"distance" binding:"omitempty,numeric,min=1"`

	Services string `form:"services"`

	Limit int64 `form:"limit" binding:"omitempty,numeric,min=1" default:"10"`
	Page  int64 `form:"page" binding:"omitempty,numeric,min=1" default:"1"`
}

type GetDressmakersByProximityOuput struct {
	*paginator.Paginate[entity.Dressmaker]
}

func (useCase *GetDressmakersUseCase) Execute(payload GetDressmakersInput) (*GetDressmakersByProximityOuput, pkg.Error) {
	fmt.Println(payload)

	var searchParams database.GetDressmakersParams

	switch {
	case payload.Services != "":
		searchParams = database.GetDressmakersParams{
			Services: payload.Services,
		}
	case payload.Latitude != 0 && payload.Longitude != 0 && payload.Distance != 0:
		searchParams = database.GetDressmakersParams{
			Latitude:  payload.Latitude,
			Longitude: payload.Longitude,
			Distance:  payload.Distance,
		}
	default:
		searchParams = database.GetDressmakersParams{
			Default: true,
		}
	}

	dressmakers, err := useCase.DressMakerRepository.Find(searchParams)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	offset := paginator.GetOffset(payload.Limit, payload.Page, dressmakers)
	paginatedItems := dressmakers[offset.Start:offset.End]

	response := &GetDressmakersByProximityOuput{
		Paginate: &paginator.Paginate[entity.Dressmaker]{
			Items:          &paginatedItems,
			PaginationInfo: paginator.NewPaginatation(payload.Limit, payload.Page, int64(len(dressmakers))),
		},
	}

	return response, pkg.Error{}
}
