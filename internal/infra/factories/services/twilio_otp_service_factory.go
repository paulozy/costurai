package factories

import (
	"github.com/paulozy/costurai/internal/infra/server/types"
	services "github.com/paulozy/costurai/internal/infra/services/otp"
)

type TwilioOTPServiceFactory struct {
	TwilioService *services.TwilioService
}

func TwilioServiceFactory(
	cfg types.TwilioConfig,
) *TwilioOTPServiceFactory {
	twilioOtpService := services.NewTwilioService(
		cfg.TwilioAccountSID,
		cfg.TwilioAuthToken,
		cfg.TwilioSMSSID,
	)

	return &TwilioOTPServiceFactory{
		TwilioService: twilioOtpService,
	}
}
