package usecases

import (
	"github.com/paulozy/costurai/internal/entity"
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/pkg"
)

type AddDressmakerReviewUseCase struct {
	DressmakerRepository        database.DressmakerRepositoryInterface
	DressmakerReviewsRepository database.DressmakerReviewsRepositoryInterface
}

type AddDressmakerReviewUseCaseInput struct {
	DressmakerID string  `json:"dressmaker_id"`
	UserID       string  `json:"user_id"`
	Comment      string  `json:"comment"`
	Grade        float64 `json:"grade"`
}

func NewAddDressmakerReviewUseCase(dmRepo database.DressmakerRepositoryInterface, dmrRepo database.DressmakerReviewsRepositoryInterface) *AddDressmakerReviewUseCase {
	return &AddDressmakerReviewUseCase{
		DressmakerRepository:        dmRepo,
		DressmakerReviewsRepository: dmrRepo,
	}
}

func (usecase *AddDressmakerReviewUseCase) Execute(input AddDressmakerReviewUseCaseInput) (*entity.Dressmaker, pkg.Error) {
	validationErr := validateNewReviewInput(input)
	if validationErr.Message != "" {
		return nil, validationErr
	}

	dressmaker, err := usecase.DressmakerRepository.FindByID(input.DressmakerID)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	if dressmaker == nil {
		return nil, pkg.NewNotFoundError("dressmaker")
	}

	review := entity.NewReview(input.DressmakerID, input.UserID, input.Comment, input.Grade)

	// dressmaker.AddReview(*review)

	err = usecase.DressmakerRepository.Update(dressmaker)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	err = usecase.DressmakerReviewsRepository.Create(review)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	return dressmaker, pkg.Error{}
}

func validateNewReviewInput(input AddDressmakerReviewUseCaseInput) pkg.Error {
	if input.DressmakerID == "" {
		return pkg.NewMissingFieldError("dressmaker_id")
	}

	if input.UserID == "" {
		return pkg.NewMissingFieldError("user_id")
	}

	if input.Comment == "" {
		return pkg.NewMissingFieldError("comment")
	}

	if input.Grade < 1 || input.Grade > 5 {
		return pkg.Error{
			Message: "grade must be between 1 and 5",
			Status:  400,
		}
	}

	return pkg.Error{}
}
