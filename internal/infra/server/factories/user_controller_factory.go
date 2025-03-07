package factories

import (
	"cloud.google.com/go/firestore"
	"github.com/paulozy/costurai/internal/infra/database/firestore/repositories"
	"github.com/paulozy/costurai/internal/infra/server/controllers"
	usecases "github.com/paulozy/costurai/internal/usecase"
)

func UserControllerFactory(
	databaseInstance *firestore.Client,
) *controllers.UserController {
	userRepository := repositories.NewFirestoreUserRepository(databaseInstance)

	createUserUseCase := usecases.NewCreateUserUseCase(userRepository)

	userUseCases := controllers.UserUseCasesInput{
		CreateUserUseCase: createUserUseCase,
	}

	userController := controllers.NewUserController(userRepository, userUseCases)

	return userController
}
