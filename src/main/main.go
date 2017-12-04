package main

import (
	"clustertruck"
	"log"
	"net/http"
	"fmt"
)

const (
	address = "0.0.0.0"
	port = 8090
)

func main() {
	httpMux := clustertruck.SetupAPI()

	log.Printf("Server running on address and port %s:%d\n", address, port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", address, port), httpMux)
	if err != nil {
		log.Fatal("Server shutdown with error: " + err.Error())
	}
}
