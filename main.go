package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/andrewburian/powermux"
	"github.com/go-pg/pg"
	"net/http"
	"os"
)

const (
	ROUTE_AUTH     = "/auth"
	ROUTE_PASSWORD = "/password"
)

const (
	CONF_PORT = "PORT"
)

func main() {

	debug := flag.Bool("debug", false, "Include verbose debug messages")
	quiet := flag.Bool("quiet", false, "Suppress all but error messages")
	flag.Parse()

	// Database connection
	authDB := &PGAuthDatabase{
		Database: pg.Connect(&pg.Options{
			Addr:     "localhost:5432", //TODO environment variables
			User:     "postgres",
			Database: "postgres",
			Password: "root",
		}),
	}

	// Create the password auth handler
	passwordHandler := PasswordAuthHandler{
		DB: authDB,
	}

	// Create the router

	// ALL routes must be under /auth
	mux := powermux.NewServeMux()

	// Setup middleware

	// debug supersedes quiet
	if *quiet {
		logrus.SetLevel(logrus.ErrorLevel)
	}
	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Running at DEBUG verbosity")
	}

	logMid := &LoggerMiddleware{
		Base: logrus.NewEntry(logrus.StandardLogger()),
	}

	mux.Route("/").Middleware(logMid)

	// Register the handlers
	authRoute := mux.Route(ROUTE_AUTH)
	passwordHandler.Setup(authRoute.Route(ROUTE_PASSWORD))

	// get the port
	var port string
	if port = os.Getenv(CONF_PORT); port == "" {
		port = "http"
	}

	// start the http server
	logrus.WithField("port", port).Info("Server starting")
	http.ListenAndServe(":"+port, mux)
}
