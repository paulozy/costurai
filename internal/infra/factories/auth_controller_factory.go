package factories

import (
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/internal/infra/database/firestore/repositories"
	"github.com/paulozy/costurai/internal/infra/server/controllers"
	"github.com/paulozy/costurai/internal/infra/server/types"
	services "github.com/paulozy/costurai/internal/infra/services/otp"
	usecases "github.com/paulozy/costurai/internal/usecase"
)

func AuthControllerFactory(
	databaseInstance database.DatabaseInterface,
	twilio types.TwilioConfig,
) *controllers.AuthController {
	dressMakerRepository := repositories.NewFirestoreDressmakerRepository(databaseInstance)
	userRepository := repositories.NewFirestoreUserRepository(databaseInstance)

	repositories := usecases.NewAuthenticationUseCaseInput{
		DressMakerRepository: dressMakerRepository,
		UserRepository:       userRepository,
	}

	twilioOtpService := services.NewTwilioService(
		twilio.TwilioAccountSID,
		twilio.TwilioAuthToken,
		twilio.TwilioSMSSID,
	)

	authenticationUseCase := usecases.NewDressMakerAuthenticationUseCase(repositories)
	sendOtpCodeUseCase := usecases.NewSendOTPUseCase(twilioOtpService)

	authUseCasesInput := controllers.AuthUseCasesInput{
		AuthUseCase:        authenticationUseCase,
		SendOtpCodeUseCase: sendOtpCodeUseCase,
	}

	authController := controllers.NewAuthController(authUseCasesInput, dressMakerRepository)

	return authController
}
