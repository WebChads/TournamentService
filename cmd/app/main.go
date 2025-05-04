package main

import (
	"context"
	"log/slog"
	"os"

	prettylogger "github.com/WebChads/AccountService/pkg/pretty_logger"
	_ "github.com/WebChads/TournamentService/docs"
	"github.com/WebChads/TournamentService/internal/config"
	server "github.com/WebChads/TournamentService/internal/delivery/http"
	slogerr "github.com/WebChads/TournamentService/internal/pkg/logger"
	"github.com/WebChads/TournamentService/internal/storage/pgsql/migrations"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:5000
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// Init config
	config, env := config.NewServerConfig()
	if config == nil || env == nil {
		return
	}

	// Init logger
	logger := setupLogger(config.LogLevel)

	// Create context
	ctx := context.Background()

	// Init database
	db, err := server.NewDB(ctx, env.DatabaseURL)
	if err != nil {
		logger.Error("failed to create database", slogerr.Error(err))
		return
	}
	defer db.Close()

	// Apply migrations
	if err := migrations.RunMigrations(db.DB, logger); err != nil {
		return
	}

	// Configure server
	router := server.InitRouter(config, logger, db)
	srv := server.NewServer(router, env.Address)

	// Run server
	logger.Info("server started", "address", env.Address)
	srv.ListenAndServe()
}

const (
	logLocal = "local"
	logStage = "stage"
	logProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case logLocal:
		handler := prettylogger.NewPrettyHandler(os.Stdout)
		log = slog.New(handler)
	case logStage:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case logProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
