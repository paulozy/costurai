package usecases

import (
	"fmt"

	services "github.com/paulozy/costurai/internal/infra/services/sms"
	"github.com/paulozy/costurai/internal/usecase/auth/dtos"
	"github.com/paulozy/costurai/pkg"
)

type SendOTPUseCase struct {
	OTPService services.OTPServiceInterface
}

type NewSendOTPUseCaseInput struct {
	OTPService services.OTPServiceInterface
}

func NewSentOTPUseCase(services NewSendOTPUseCaseInput) *SendOTPUseCase {
	return &SendOTPUseCase{
		OTPService: services.OTPService,
	}
}

func (useCase *SendOTPUseCase) Execute(payload dtos.SendOTPInput) pkg.Error {
	err := useCase.OTPService.Send(payload.Phone)
	if err != nil {
		fmt.Println(err)
		return pkg.Error{
			Error:   err.Error(),
			Message: "Error on sending code",
		}
	}

	return pkg.Error{}
}
