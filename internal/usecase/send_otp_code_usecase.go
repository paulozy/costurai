package usecases

import (
	services "github.com/paulozy/costurai/internal/infra/services/otp"
	"github.com/paulozy/costurai/pkg"
)

type SendOTPCodeUseCase struct {
	OTPService services.OTPServiceInterface
}

func NewSendOTPUseCase(otp services.OTPServiceInterface) *SendOTPCodeUseCase {
	return &SendOTPCodeUseCase{
		OTPService: otp,
	}
}

type SendOTPCodeInput struct {
	Phone string `json:"phone"`
}

func (useCase *SendOTPCodeUseCase) Execute(input SendOTPCodeInput) (*services.OTPServiceSendVerificationOutput, pkg.Error) {
	response, err := useCase.OTPService.SendVerification(input.Phone)
	if err != nil {
		return nil, pkg.Error{
			Error:   err.Error(),
			Message: "Error on send OTP code",
		}
	}

	return response, pkg.Error{}
}
