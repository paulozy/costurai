package main

import (
	"fmt"
	"os"

	"github.com/paulozy/costurai/configs"
	"github.com/paulozy/costurai/internal/infra/database/firestore"
	"github.com/paulozy/costurai/internal/infra/server"
)

const (
	Firestore string = "firestore"
	Postgres  string = "postgres"
)

func main() {
	fmt.Println("Starting the Costurai API server...")

	env := "development"
	if len(os.Args) > 1 {
		env = os.Args[1]
	}
	if env != "local" && env != "development" && env != "production" {
		panic(fmt.Sprintf("Invalid environment: %s. Allowed values are 'local', 'development', or 'production'.", env))
	}

	configs, err := configs.LoadConfig(env)
	if err != nil {
		panic(err)
	}

	database := firestore.NewFirestoreClient(configs.FirebaseProjectId)

	server := server.NewServer(configs.WebHost, configs.WebPort, configs.Env)
	server.FirestoreDB = database
	server.AddHandlers()
	server.Start()
}
