package main

import (
	"log"
	"net/http"

	"github.com/rkcpi/vell/api"
	"github.com/rkcpi/vell/config"
)

func main() {
	EnsureCreateRepo()

	router := api.NewRouter()

	log.Printf("Vell repositories location: %s", config.ReposPath)

	EnsureWritableReposPath()

	log.Printf("Listening for requests on %s", config.ListenAddress)
	log.Fatal(http.ListenAndServe(config.ListenAddress, router))
}
