package main

import (
	"context"
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/andrewburian/powermux"
	"github.com/go-pg/pg"
	"github.com/kelseyhightower/envconfig"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	// We allow all methods through CORS, the router will respond to
	// verbs that aren't implemented
	corsMid := &CorsHandler{
		AllowedOrigins: []string{conf.WebappOrigin},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}

	mux.Route(ROUTE_AUTH).Middleware(corsMid).Options(corsMid)

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

	// setup the http server
	logrus.WithField("port", conf.Port).Info("Server starting")
	server := &http.Server{
		Addr:    ":" + conf.Port,
		Handler: mux,
	}

	// Trap TERM and INT signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// Signals kill the server
	go func(c <-chan os.Signal) {
		select {
		case sig := <-c:
			shutdownTime := time.Duration(conf.ShutdownTime) * time.Second
			shutdownCtx, cancelFunc := context.WithTimeout(context.Background(), shutdownTime)
			logrus.WithField("signal", sig).Warn("Trapped signal")
			server.Shutdown(shutdownCtx)
			cancelFunc()
		}
	}(sigChan)

	// Run the server
	err = server.ListenAndServe()

	// clean exit on server close
	if err == http.ErrServerClosed {
		logrus.Info("Server shut down")
		return
	}

	// Error otherwise
	logrus.Fatal(err)
}
