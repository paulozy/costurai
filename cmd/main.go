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

	var env *string
	if len(os.Args) > 1 {
		env = &os.Args[1]
	} else {
		env = new(string)
		*env = "local"
	}

	configs, err := configs.LoadConfig(*env)
	if err != nil {
		panic(err)
	}

	database := firestore.NewFirestoreClient(configs.FirebaseProjectId)

	server := server.NewServer(configs.WebHost, configs.WebPort, configs.Env)
	server.FirestoreDB = database
	server.AddHandlers()
	server.Start()
}
