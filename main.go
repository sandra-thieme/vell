package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rkcpi/vell/routing"
)

func main() {

	router := routing.NewRouter()

	var port string
	if port = os.Getenv("VELL_HTTP_PORT"); port == "" {
		port = "8080"
	}

	var address string
	if address = os.Getenv("VELL_HTTP_ADDRESS"); address == "" {

	}
	listenAddress := fmt.Sprintf("%s:%s", address, port)

	log.Fatal(http.ListenAndServe(listenAddress, router))
}
