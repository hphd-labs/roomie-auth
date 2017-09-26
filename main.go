package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/andrewburian/powermux"
	"github.com/go-pg/pg"
	"github.com/kelseyhightower/envconfig"
	"net/http"
)

const (
	ROUTE_AUTH     = "/auth"
	ROUTE_PASSWORD = "/password"
)

func main() {

	debug := flag.Bool("debug", false, "Include verbose debug messages")
	quiet := flag.Bool("quiet", false, "Suppress all but error messages")
	flag.Parse()

	// Parse environment variables
	var conf Config
	if err := envconfig.Process("", &conf); err != nil {
		logrus.Fatal(err)
		return
	}

	// Database connection
	dbOpts, err := pg.ParseURL(conf.DatabaseUrl)
	if err != nil {
		logrus.Fatal(err)
		return
	}

	authDB := &PGAuthDatabase{
		Database: pg.Connect(dbOpts),
	}

	// Create the password auth handler
	passwordHandler := PasswordAuthHandler{
		DB:           authDB,
		PasswordCost: conf.PasswordCost,
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

	// start the http server
	logrus.WithField("port", conf.Port).Info("Server starting")
	err = http.ListenAndServe(":"+conf.Port, mux)
	logrus.Fatal(err)
}
