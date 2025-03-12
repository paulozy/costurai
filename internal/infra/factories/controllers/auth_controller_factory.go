package factories

import (
	"github.com/paulozy/costurai/internal/infra/database"
	factoryRepositories "github.com/paulozy/costurai/internal/infra/factories/repositories"
	factoryServices "github.com/paulozy/costurai/internal/infra/factories/services"
	"github.com/paulozy/costurai/internal/infra/server/controllers"
	"github.com/paulozy/costurai/internal/infra/server/types"
	usecases "github.com/paulozy/costurai/internal/usecase"
)

func AuthControllerFactory(
	databaseInstance database.DatabaseInterface,
	twilio types.TwilioConfig,
) *controllers.AuthController {
	repositories := factoryRepositories.FirestoreRepositoriesFactory(databaseInstance)

	repositoriesInput := usecases.NewAuthenticationUseCaseInput{
		DressMakerRepository: repositories.DressmakerRepository,
		UserRepository:       repositories.UserRepository,
	}
	otpServices := factoryServices.TwilioServiceFactory(twilio)

	authenticationUseCase := usecases.NewDressMakerAuthenticationUseCase(repositoriesInput)
	sendOtpCodeUseCase := usecases.NewSendOTPUseCase(otpServices.TwilioService)

	authUseCasesInput := controllers.AuthUseCasesInput{
		AuthUseCase:        authenticationUseCase,
		SendOtpCodeUseCase: sendOtpCodeUseCase,
	}

	authController := controllers.NewAuthController(authUseCasesInput)

	return authController
}
