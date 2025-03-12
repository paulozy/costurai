package factories

import (
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/internal/infra/database/firestore/repositories"
	"github.com/paulozy/costurai/internal/infra/server/controllers"
	"github.com/paulozy/costurai/internal/infra/server/types"
	services "github.com/paulozy/costurai/internal/infra/services/otp"
	usecases "github.com/paulozy/costurai/internal/usecase"
)

func DressmakerControllerFactory(
	databaseInstance database.DatabaseInterface,
	twilio types.TwilioConfig,
) *controllers.DressmakerController {

	dressMakerRepository := repositories.NewFirestoreDressmakerRepository(databaseInstance)
	dressMakerReviewsRepository := repositories.NewFirestoreReviewsRepository(databaseInstance)

	twilioOtpService := services.NewTwilioService(
		twilio.TwilioAccountSID,
		twilio.TwilioAuthToken,
		twilio.TwilioSMSSID,
	)

	createDressmakerUseCase := usecases.NewCreateDressMakerUseCase(dressMakerRepository)
	showDressmakerUseCase := usecases.NewShowDressmakerUseCase(dressMakerRepository)
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
		ShowDressmakerUseCase:            showDressmakerUseCase,
	}

	dressmakerController := controllers.NewDressmakerController(dressMakerRepository, dressMakerReviewsRepository, dressmakerUseCases)

	return dressmakerController
}
