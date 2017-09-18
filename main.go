package main

import (
	"github.com/andrewburian/powermux"
	"net/http"
	"os"
)

const (
	ROUTE_AUTH = "/auth"
)

const (
	CONF_PORT = "PORT"
)

func main() {

	// Create the password auth handler
	passwordHandler := PasswordAuthHandler{}

	// Create the router
	// ALL routes must be unde
	mux := powermux.NewServeMux()

	// Register the handlers
	authRoute := mux.Route(ROUTE_AUTH)
	passwordHandler.Setup(authRoute.Route("/password"))

	// get the port
	var port string
	if port = os.Getenv(CONF_PORT); port == "" {
		port = "http"
	}

	// start the http server
	http.ListenAndServe(":"+port, mux)
}
