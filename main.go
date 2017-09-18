package main

import (
	"github.com/andrewburian/powermux"
	"net/http"
	"os"
)

const (
	ROUTE_AUTH = "/auth"
	ROUTE_PASSWORD = "/password"
)

const (
	CONF_PORT = "PORT"
)

func main() {

	// Database connection
	authDB := &AuthDatabase{}

	// Create the password auth handler
	passwordHandler := PasswordAuthHandler{
		DB: authDB,
	}

	// Create the router
	// ALL routes must be unde
	mux := powermux.NewServeMux()

	// Register the handlers
	authRoute := mux.Route(ROUTE_AUTH)
	passwordHandler.Setup(authRoute.Route(ROUTE_PASSWORD))

	// get the port
	var port string
	if port = os.Getenv(CONF_PORT); port == "" {
		port = "http"
	}

	// start the http server
	http.ListenAndServe(":"+port, mux)
}
