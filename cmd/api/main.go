package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const appVersion = "1.0.0"

type serverConfig struct {
	port        int
	environment string
}

type applicationDependecies struct {
	config serverConfig
	logger *slog.Logger
}

func main() {

	var settings serverConfig

	flag.IntVar(&settings.port, "port", 4000, "Server port")
	flag.StringVar(&settings.environment, "env", "development", "Environment(developmnet|staging|production)")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	appInstance := &applicationDependecies{
		config: settings,
		logger: logger,
	}

	router := http.NewServeMux()
	router.HandleFunc("/v1/healthcheck", appInstance.healthChechHandler)

	apiServer := &http.Server{
		Addr:        fmt.Sprintf(":%d", settings.port),
		Handler:     router,
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		ErrorLog:    slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "address", apiServer.Addr, "environment", settings.environment)
	err := apiServer.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)

}
