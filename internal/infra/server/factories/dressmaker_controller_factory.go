package factories

import (
	"cloud.google.com/go/firestore"
	"github.com/paulozy/costurai/internal/infra/database/firestore/repositories"
	"github.com/paulozy/costurai/internal/infra/server/controllers"
	sInterfaces "github.com/paulozy/costurai/internal/infra/server/interfaces"
	services "github.com/paulozy/costurai/internal/infra/services/otp"
	usecases "github.com/paulozy/costurai/internal/usecase"
)

func DressmakerControllerFactory(
	databaseInstance *firestore.Client,
	twilio sInterfaces.TwilioConfig,
) *controllers.DressmakerController {

	dressMakerRepository := repositories.NewFirestoreDressmakerRepository(databaseInstance)
	dressMakerReviewsRepository := repositories.NewFirestoreReviewsRepository(databaseInstance)

	twilioOtpService := services.NewTwilioService(
		twilio.TwilioAccountSID,
		twilio.TwilioAuthToken,
		twilio.TwilioSMSSID,
	)

	createDressmakerUseCase := usecases.NewCreateDressMakerUseCase(dressMakerRepository)
	enableDresskamerUseCase := usecases.NewEnableDressmakerUseCase(dressMakerRepository, twilioOtpService)
	updateDressmakerUseCase := usecases.NewUpdateDressMakerUseCase(dressMakerRepository)
	getDressmakersByProximityUseCase := usecases.NewGetDressmakersByProximityUseCase(dressMakerRepository)
	addDressmakerReviewUseCase := usecases.NewAddDressmakerReviewUseCase(dressMakerRepository, dressMakerReviewsRepository)

	dressmakerUseCases := controllers.DressmakerUseCasesInput{
		CreateDressmakerUseCase:          createDressmakerUseCase,
		UpdateDressmakerUseCase:          updateDressmakerUseCase,
		GetDressmakersByProximityUseCase: getDressmakersByProximityUseCase,
		AddDressmakerReviewUseCase:       addDressmakerReviewUseCase,
		EnableDressmakerUseCase:          enableDresskamerUseCase,
	}

	dressmakerController := controllers.NewDressmakerController(dressMakerRepository, dressMakerReviewsRepository, dressmakerUseCases)

	return dressmakerController
}
