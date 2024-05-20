package server

import (
	"database/sql"

	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/internal/infra/server/controllers"
	usecases "github.com/paulozy/costurai/internal/usecase"
)

var Routes = []Handler{}

func PopulateRoutes(db *sql.DB) []Handler {
	addDressmakerRoutes(db)
	addUserRoutes(db)
	addAuthRoutes(db)
	return Routes
}

func addDressmakerRoutes(db *sql.DB) {
	dressMakerRepository := database.NewDressmakerRepository(db)
	dressMakerReviewsRepository := database.NewDressmakerReviewsRepository(db)
	createDressmakerUseCase := usecases.NewCreateDressMakerUseCase(dressMakerRepository)
	updateDressmakerUseCase := usecases.NewUpdateDressMakerUseCase(dressMakerRepository)
	getDressmakersByProximityUseCase := usecases.NewGetDressmakersByProximityUseCase(dressMakerRepository)
	getDressmakersByServicesUseCase := usecases.NewGetDressmakersByServicesUseCase(dressMakerRepository)
	addDressmakerReviewUseCase := usecases.NewAddDressmakerReviewUseCase(dressMakerRepository, dressMakerReviewsRepository)

	dressmakerUseCases := controllers.DressmakerUseCasesInput{
		CreateDressmakerUseCase: createDressmakerUseCase,
		UpdateDressmakerUseCase: updateDressmakerUseCase,
		GetDressmakersByProximityUseCase: getDressmakersByProximityUseCase,
		GetDressmakersByServicesUseCase: getDressmakersByServicesUseCase,
		AddDressmakerReviewUseCase: addDressmakerReviewUseCase,
	}

	dressmakerController := controllers.NewDressmakerController(dressMakerRepository, dressMakerReviewsRepository, dressmakerUseCases)

	dressmakerControllerRoutes := []Handler{
		{
			Path:   "/dressmakers",
			Method: "POST",
			Func:   dressmakerController.CreateDressmaker,
		},
		{
			Path:   "/dressmakers",
			Method: "GET",
			Func:   dressmakerController.GetDressmakers,
		},
		{
			Path:   "/dressmakers/:id",
			Method: "PUT",
			Auth:   true,
			Func:   dressmakerController.UpdateDressmaker,
		},
		{
			Path:   "/dressmakers/reviews",
			Method: "POST",
			Auth:   true,
			Func:  dressmakerController.AddDressmakerReview,
		},
	}

	Routes = append(Routes, dressmakerControllerRoutes...)
}

func addUserRoutes(db *sql.DB) {
	userRepository := database.NewUserRepository(db)
	createUserUseCase := usecases.NewCreateUserUseCase(userRepository)

	userUseCases := controllers.UserUseCasesInput{
		CreateUserUseCase: createUserUseCase,
	}

	userController := controllers.NewUserController(userRepository, userUseCases)

	userControllerRoutes := []Handler{
		{
			Path:   "/users",
			Method: "POST",
			Func:   userController.CreateUser,
		},
	}

	Routes = append(Routes, userControllerRoutes...)
}

func addAuthRoutes(db *sql.DB) {
	dressMakerRepository := database.NewDressmakerRepository(db)
	userRepository := database.NewUserRepository(db)

	repositories := usecases.NewAuthenticationUseCaseInput{
		DressMakerRepository: dressMakerRepository,
		UserRepository:       userRepository,
	}

	authenticationUseCase := usecases.NewDressMakerAuthenticationUseCase(repositories)

	authController := controllers.NewAuthController(authenticationUseCase, dressMakerRepository)

	dressmakerAuthHandler := Handler{
		Path:   "/dressmakers/auth",
		Method: "POST",
		Func:   authController.AuthenticateDressmaker,
	}

	userAuthHandler := Handler{
		Path:   "/users/auth",
		Method: "POST",
		Func:   authController.AuthenticateUser,
	}

	Routes = append(Routes, dressmakerAuthHandler, userAuthHandler)
}

