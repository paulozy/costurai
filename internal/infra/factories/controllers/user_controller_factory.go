package factories

import (
	"github.com/paulozy/costurai/internal/infra/database"
	factoryRepositories "github.com/paulozy/costurai/internal/infra/factories/repositories"
	"github.com/paulozy/costurai/internal/infra/server/controllers"
	usecases "github.com/paulozy/costurai/internal/usecase"
)

func UserControllerFactory(
	databaseInstance database.DatabaseInterface,
) *controllers.UserController {
	repositories := factoryRepositories.FirestoreRepositoriesFactory(databaseInstance)

	createUserUseCase := usecases.NewCreateUserUseCase(repositories.UserRepository)

	userUseCases := controllers.UserUseCasesInput{
		CreateUserUseCase: createUserUseCase,
	}

	userController := controllers.NewUserController(userUseCases)

	return userController
}
