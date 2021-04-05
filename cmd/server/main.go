package main

import (
	"database/sql"
	"fmt"
	"os"
	"weight-tracker/pkg/api"
	"weight-tracker/pkg/app"
	"weight-tracker/pkg/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "this is the startup error: %s\\n", err)
		os.Exit(1)
	}
}

// func run will be responsible for setting up db connections, routers etc
func run() error {
	connectionString := "postgres://postgres:postgres@localhost:5432/weight_tracker?sslmode=disable"

	db, err := setupDatabase(connectionString)

	if err != nil {
		return err
	}

	// create storage dependency
	storage := repository.NewStorage(db)

	// create router dependecy
	router := gin.Default()
	router.Use(cors.Default())

	// create user service
	userService := api.NewUserService(storage)

	server := app.NewServer(router, userService)

	// start the server
	err = server.Run()

	if err != nil {
		return err
	}

	return nil
}

// func setupDatabase setup database connection
func setupDatabase(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)

	if err != nil {
		return nil, err
	}

	// ping the DB to ensure that it is connected
	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil

}
