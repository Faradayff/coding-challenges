package main

import (
	"log"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
)

const (
	ListenAddress = ":8080"
	// TODO: add further configuration parameters here ...
)

// @title Signing Service API
// @version 1.0
// @description API for managing signature devices and signing transactions.
// @host http://localhost:8080/
// @BasePath /api/v0

func main() {
	server := api.NewServer(ListenAddress)

	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
