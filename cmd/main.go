package main

import (
	"fmt"

	"github.com/paulozy/costurai/configs"
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/internal/infra/server"
)

func main() {
	fmt.Println("Starting the Costurai API server...")

	configs, err := configs.LoadConfig("../")
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", configs.DBHost, configs.DBPort, configs.DBUser, configs.DBPassword, configs.DBName)
	db := database.NewDatabaseConn(dsn)

	server.PopulateRoutes(db)

	server := server.NewServer(configs.WebHost, configs.WebPort)
	server.AddHandlers()
	server.Start()
}
