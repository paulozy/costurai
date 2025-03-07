package factories

import (
	"cloud.google.com/go/firestore"
	"github.com/paulozy/costurai/internal/infra/database/firestore/repositories"
	"github.com/paulozy/costurai/internal/infra/server/controllers"
	usecases "github.com/paulozy/costurai/internal/usecase"
)

func AuthControllerFactory(
	databaseInstance *firestore.Client,
) *controllers.AuthController {
	dressMakerRepository := repositories.NewFirestoreDressmakerRepository(databaseInstance)
	userRepository := repositories.NewFirestoreUserRepository(databaseInstance)

	repositories := usecases.NewAuthenticationUseCaseInput{
		DressMakerRepository: dressMakerRepository,
		UserRepository:       userRepository,
	}

	authenticationUseCase := usecases.NewDressMakerAuthenticationUseCase(repositories)

	authController := controllers.NewAuthController(authenticationUseCase, dressMakerRepository)

	return authController
}
