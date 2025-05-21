package server

import (
	"cloud.google.com/go/firestore"
	"github.com/paulozy/costurai/internal/infra/database/firestore/repositories"
	"github.com/paulozy/costurai/internal/infra/server/controllers"
	services "github.com/paulozy/costurai/internal/infra/services/sms"
	authUseCases "github.com/paulozy/costurai/internal/usecase/auth"
	dressmakerUseCases "github.com/paulozy/costurai/internal/usecase/dressmaker"
	userUseCases "github.com/paulozy/costurai/internal/usecase/user"
)

var Routes = []Handler{}

func PopulateRoutes(db *firestore.Client) []Handler {
	addDressmakerRoutes(db)
	addUserRoutes(db)
	addAuthRoutes(db)
	return Routes
}

func addDressmakerRoutes(db *firestore.Client) {
	dressMakerRepository := repositories.NewFirestoreDressmakerRepository(db)

	createDressmakerUseCase := dressmakerUseCases.NewCreateDressMakerUseCase(dressMakerRepository)
	updateDressmakerUseCase := dressmakerUseCases.NewUpdateDressMakerUseCase(dressMakerRepository)
	getDressmakersByProximityUseCase := dressmakerUseCases.NewGetDressmakersByProximityUseCase(dressMakerRepository)
	showDressmakerUseCase := dressmakerUseCases.NewShowDressMakerUseCase(dressMakerRepository)

	dressmakerUseCases := controllers.DressmakerUseCasesInput{
		CreateDressmakerUseCase:          createDressmakerUseCase,
		UpdateDressmakerUseCase:          updateDressmakerUseCase,
		GetDressmakersByProximityUseCase: getDressmakersByProximityUseCase,
		ShowDressmakerUseCase:            showDressmakerUseCase,
	}

	dressmakerController := controllers.NewDressmakerController(dressMakerRepository, nil, dressmakerUseCases)

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
			Method: "GET",
			Func:   dressmakerController.GetDressmaker,
		},
		{
			Path:   "/dressmakers/:id",
			Method: "PUT",
			Auth:   true,
			Func:   dressmakerController.UpdateDressmaker,
		},
	}

	Routes = append(Routes, dressmakerControllerRoutes...)
}

func addUserRoutes(db *firestore.Client) {
	userRepository := repositories.NewFirestoreUserRepository(db)

	createUserUseCase := userUseCases.NewCreateUserUseCase(userRepository)

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

func addAuthRoutes(db *firestore.Client) {
	dressmakerRepository := repositories.NewFirestoreDressmakerRepository(db)
	userRepository := repositories.NewFirestoreUserRepository(db)

	authDressmakerUseCase := authUseCases.NewDressmakerAuthenticationUseCase(authUseCases.NewAuthDressmakerUseCaseInput{
		DressmakerRepository: dressmakerRepository,
	})
	authUserUseCase := authUseCases.NewUserAuthUseCase(authUseCases.NewAuthUserUseCaseInput{
		UserRepository: userRepository,
	})

	OTPService := services.NewTwilioService()
	sendOTPUseCase := authUseCases.NewSentOTPUseCase(
		authUseCases.NewSendOTPUseCaseInput{
			OTPService: OTPService,
		},
	)
	verifyOTPUseCase := authUseCases.NewVerifyOTPUseCase(authUseCases.NewVerifyOTPUseCaseInput{
		OTPService:           OTPService,
		DressmakerRepository: dressmakerRepository,
		UserRepository:       userRepository,
	})

	authController := controllers.NewAuthController(
		authDressmakerUseCase,
		authUserUseCase,
		sendOTPUseCase,
		verifyOTPUseCase,
		dressmakerRepository,
		userRepository,
	)

	authHandlers := []Handler{
		{
			Path:   "/dressmakers/auth",
			Method: "POST",
			Func:   authController.AuthenticateDressmaker,
		},
		{
			Path:   "/users/auth",
			Method: "POST",
			Func:   authController.AuthenticateUser,
		},
		{
			Path:   "/otp",
			Method: "POST",
			Func:   authController.SendOTP,
		},
		{
			Path:   "/otp/dressmaker/verify",
			Method: "POST",
			Func:   authController.VerifyOTP,
			Auth:   true,
		},
		{
			Path:   "/otp/user/verify",
			Method: "POST",
			Func:   authController.VerifyOTP,
			Auth:   true,
		},
	}

	Routes = append(Routes, authHandlers...)
}
