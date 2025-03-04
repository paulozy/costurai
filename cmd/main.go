package main

import (
	"fmt"

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

	configs, err := configs.LoadConfig("../")
	if err != nil {
		panic(err)
	}

	database := firestore.NewFirestoreClient(configs.FirebaseProjectId)

	server := server.NewServer(configs.WebHost, configs.WebPort, configs.Env)
	server.FirestoreDB = database
	server.AddHandlers()
	server.Start()
}
