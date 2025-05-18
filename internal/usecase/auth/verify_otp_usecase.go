package usecases

import (
	"github.com/paulozy/costurai/internal/infra/database"
	services "github.com/paulozy/costurai/internal/infra/services/sms"
	"github.com/paulozy/costurai/internal/usecase/auth/dtos"
	"github.com/paulozy/costurai/pkg"
)

type VerifyOTPUseCase struct {
	OTPService           services.OTPServiceInterface
	DressmakerRepository database.DressmakerRepositoryInterface
	UserRepository       database.UserRepositoryInterface
}

type NewVerifyOTPUseCaseInput struct {
	OTPService           services.OTPServiceInterface
	DressmakerRepository database.DressmakerRepositoryInterface
	UserRepository       database.UserRepositoryInterface
}

func NewVerifyOTPUseCase(input NewVerifyOTPUseCaseInput) *VerifyOTPUseCase {
	return &VerifyOTPUseCase{
		OTPService:           input.OTPService,
		DressmakerRepository: input.DressmakerRepository,
		UserRepository:       input.UserRepository,
	}
}

func (uc *VerifyOTPUseCase) Execute(payload dtos.VerifyOTPInput) pkg.Error {
	ok, err := uc.OTPService.Verify(payload.Phone, payload.Code)
	if err != nil {
		return pkg.Error{
			Error:   err.Error(),
			Message: "Error on verify code",
			Status:  500,
		}
	} else if !ok {
		return pkg.Error{
			Error:   "Error on verify code",
			Message: "Invalid code",
			Status:  400,
		}
	}

	switch payload.Enabling {
	case "dressmaker":
		return uc.enableDressmaker(payload.DressmakerID)
	case "user":
		return uc.enableUser(payload.UserID)
	default:
		return pkg.Error{
			Error:   "Error on verify code",
			Message: "Has no enabling entity",
			Status:  400,
		}
	}
}

func (uc *VerifyOTPUseCase) enableDressmaker(id string) pkg.Error {
	dressmaker, err := uc.DressmakerRepository.FindByID(id)
	if err != nil {
		return pkg.Error{
			Message: err.Error(),
			Error:   "Dressmaker not found",
			Status:  404,
		}
	}

	dressmaker.Enable()

	err = uc.DressmakerRepository.Update(dressmaker)
	if err != nil {
		return pkg.Error{
			Message: err.Error(),
			Error:   "Error on update dressmaker",
			Status:  500,
		}
	}

	return pkg.Error{}
}

func (uc *VerifyOTPUseCase) enableUser(id string) pkg.Error {
	user, err := uc.UserRepository.FindByID(id)
	if err != nil {
		return pkg.Error{
			Message: err.Error(),
			Error:   "User not found",
			Status:  404,
		}
	}

	user.Enable()

	err = uc.UserRepository.Update(user)
	if err != nil {
		return pkg.Error{
			Message: err.Error(),
			Error:   "Error on update user",
			Status:  500,
		}
	}

	return pkg.Error{}
}
