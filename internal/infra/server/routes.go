package server

import (
	"cloud.google.com/go/firestore"
	"github.com/paulozy/costurai/internal/infra/server/factories"
	sInterfaces "github.com/paulozy/costurai/internal/infra/server/interfaces"
)

type PopulateRoutesInput struct {
	db           *firestore.Client
	twilioConfig sInterfaces.TwilioConfig
}

var Routes = []Handler{}

func PopulateRoutes(input PopulateRoutesInput) []Handler {
	addDressmakerRoutes(input.db, input.twilioConfig)
	addUserRoutes(input.db)
	addAuthRoutes(input.db)
	return Routes
}

func addDressmakerRoutes(db *firestore.Client, twilioCfg sInterfaces.TwilioConfig) {
	dressmakerController := factories.DressmakerControllerFactory(db, twilioCfg)

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
			Func:   dressmakerController.AddDressmakerReview,
		},
		{
			Path:   "/dressmakers/:id/enable",
			Method: "PUT",
			Auth:   false,
			Func:   dressmakerController.EnableDressmaker,
		},
	}

	Routes = append(Routes, dressmakerControllerRoutes...)
}

func addUserRoutes(db *firestore.Client) {
	userController := factories.UserControllerFactory(db)

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
	authController := factories.AuthControllerFactory(db)

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
