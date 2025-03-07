package usecases

import (
	"github.com/paulozy/costurai/internal/entity"
	"github.com/paulozy/costurai/internal/infra/database"
	services "github.com/paulozy/costurai/internal/infra/services/otp"
	"github.com/paulozy/costurai/pkg"
)

type EnableDressmakerUseCase struct {
	DressMakerRepository database.DressmakerRepositoryInterface
	OTPService           services.OTPServiceInterface
}

func NewEnableDressmakerUseCase(
	repo database.DressmakerRepositoryInterface,
	otp services.OTPServiceInterface,
) *EnableDressmakerUseCase {
	return &EnableDressmakerUseCase{
		DressMakerRepository: repo,
		OTPService:           otp,
	}
}

type EnableDressmakerInput struct {
	ID    string `json:"id"`
	Code  string `json:"code"`
	Phone string `json:"phone"`
}

const (
	APPROVED = "approved"
)

func (useCase *EnableDressmakerUseCase) Execute(input EnableDressmakerInput) (*entity.Dressmaker, pkg.Error) {
	dressmaker, err := useCase.DressMakerRepository.FindByID(input.ID)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	verify, err := useCase.OTPService.Verify(input.Code, input.Phone)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	if verify.Status != APPROVED {
		return nil, pkg.NewInvalidOTPCode()
	}

	dressmaker.Enable()

	err = useCase.DressMakerRepository.Update(dressmaker)
	if err != nil {
		return nil, pkg.NewInternalServerError(err)
	}

	return dressmaker, pkg.Error{}
}
