package main

import (
	"os"

	"github.com/ratludu/grpc-habits-tracker/internal/database"
	"github.com/ratludu/grpc-habits-tracker/internal/server"
	"github.com/ratludu/grpc-habits-tracker/log"
)

const port = 8080

func main() {

	lgr := log.New(os.Stdout)
	db, err := database.Open("/tmp/habits")
	if err != nil {
		lgr.Logf("Cannot connect to database")
		os.Exit(1)
	}

	srv := server.New(lgr, db)

	err = srv.ListenAndServe(port)
	if err != nil {
		lgr.Logf("Error while running the server: %s", err.Error())
		os.Exit(1)
	}

}
