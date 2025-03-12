package factories

import (
	"github.com/paulozy/costurai/internal/infra/database"
	factoryRepositories "github.com/paulozy/costurai/internal/infra/factories/repositories"
	factoryServices "github.com/paulozy/costurai/internal/infra/factories/services"
	"github.com/paulozy/costurai/internal/infra/server/controllers"
	"github.com/paulozy/costurai/internal/infra/server/types"
	usecases "github.com/paulozy/costurai/internal/usecase"
)

func DressmakerControllerFactory(
	databaseInstance database.DatabaseInterface,
	twilio types.TwilioConfig,
) *controllers.DressmakerController {
	repositories := factoryRepositories.FirestoreRepositoriesFactory(databaseInstance)
	otpServices := factoryServices.TwilioServiceFactory(twilio)

	createDressmakerUseCase := usecases.NewCreateDressMakerUseCase(repositories.DressmakerRepository)
	showDressmakerUseCase := usecases.NewShowDressmakerUseCase(repositories.DressmakerRepository)
	enableDresskamerUseCase := usecases.NewEnableDressmakerUseCase(repositories.DressmakerRepository, otpServices.TwilioService)
	updateDressmakerUseCase := usecases.NewUpdateDressMakerUseCase(repositories.DressmakerRepository)
	getDressmakersByProximityUseCase := usecases.NewGetDressmakersByProximityUseCase(repositories.DressmakerRepository)
	addDressmakerReviewUseCase := usecases.NewAddDressmakerReviewUseCase(repositories.DressmakerRepository, repositories.DressmakerReviewRepository)

	dressmakerUseCases := controllers.DressmakerUseCasesInput{
		CreateDressmakerUseCase:          createDressmakerUseCase,
		UpdateDressmakerUseCase:          updateDressmakerUseCase,
		GetDressmakersByProximityUseCase: getDressmakersByProximityUseCase,
		AddDressmakerReviewUseCase:       addDressmakerReviewUseCase,
		EnableDressmakerUseCase:          enableDresskamerUseCase,
		ShowDressmakerUseCase:            showDressmakerUseCase,
	}

	dressmakerController := controllers.NewDressmakerController(dressmakerUseCases)

	return dressmakerController
}
