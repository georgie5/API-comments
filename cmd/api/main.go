package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"sync"
	"time"

	"georgie5.net/API-comments/internal/data"
	"georgie5.net/API-comments/internal/mailer"
	_ "github.com/lib/pq"
)

const appVersion = "1.0.0"

type serverConfig struct {
	port        int
	environment string
	db          struct {
		dsn string
	}

	limiter struct {
		rps     float64 // requests per second
		burst   int     // initial requests possible
		enabled bool    // enable or disable rate limiter
	}

	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type applicationDependencies struct {
	config       serverConfig
	logger       *slog.Logger
	commentModel data.CommentModel
	userModel    data.UserModel
	mailer       mailer.Mailer
	wg           sync.WaitGroup // need this later for background jobs

}

func main() {

	var settings serverConfig

	flag.IntVar(&settings.port, "port", 4000, "Server port")
	flag.StringVar(&settings.environment, "env", "development", "Environment(developmnet|staging|production)")

	flag.StringVar(&settings.db.dsn, "db-dsn", "postgres://comments:fishsticks@localhost/comments?sslmode=disable", "PostgreSQL DSN")

	flag.Float64Var(&settings.limiter.rps, "limiter-rps", 2, "Rate Limiter maximum requests per second")

	flag.IntVar(&settings.limiter.burst, "limiter-burst", 5, "Rate Limiter maximum burst")

	flag.BoolVar(&settings.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&settings.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	// We have port 25, 465, 587, 2525. If 25 doesn't work choose another
	flag.IntVar(&settings.smtp.port, "smtp-port", 2525, "SMTP port")
	// Use your Username value provided by Mailtrap
	flag.StringVar(&settings.smtp.username, "smtp-username", "5f68753fd3d8ce", "SMTP username")

	// Use your Password value provided by Mailtrap
	flag.StringVar(&settings.smtp.password, "smtp-password", "f8e72801757fc0", "SMTP password")

	flag.StringVar(&settings.smtp.sender, "smtp-sender", "Comments Community <no-reply@commentscommunity.georgie.net>", "SMTP sender")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// the call to openDB() sets up our connection pool
	db, err := openDB(settings)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// release the database resources before exiting
	defer db.Close()

	logger.Info("database connection pool established")

	appInstance := &applicationDependencies{
		config:       settings,
		logger:       logger,
		commentModel: data.CommentModel{DB: db},
		userModel:    data.UserModel{DB: db},
		mailer:       mailer.New(settings.smtp.host, settings.smtp.port, settings.smtp.username, settings.smtp.password, settings.smtp.sender),
	}

	err = appInstance.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}

func openDB(settings serverConfig) (*sql.DB, error) {
	// open a connection pool
	db, err := sql.Open("postgres", settings.db.dsn)
	if err != nil {
		return nil, err
	}

	// set a context to ensure DB operations don't take too long
	ctx, cancel := context.WithTimeout(context.Background(),
		5*time.Second)
	defer cancel()

	// let's test if the connection pool was created
	// we trying pinging it with a 5-second timeout
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	// return the connection pool (sql.DB)
	return db, nil

}
