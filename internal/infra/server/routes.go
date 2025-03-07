package server

import (
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/internal/infra/factories"
	"github.com/paulozy/costurai/internal/infra/server/types"
)

var Routes = []Handler{}

func PopulateRoutes(server Server) []Handler {
	addDressmakerRoutes(
		server.DatabaseInstance,
		server.Twilio,
	)
	addUserRoutes(server.DatabaseInstance)
	addAuthRoutes(server.DatabaseInstance)
	return Routes
}

func addDressmakerRoutes(
	db database.DatabaseInterface,
	twilio types.TwilioConfig,
) {
	dressmakerController := factories.DressmakerControllerFactory(
		db,
		twilio,
	)

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

func addUserRoutes(db database.DatabaseInterface) {
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

func addAuthRoutes(db database.DatabaseInterface) {
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
