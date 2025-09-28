package main

import (
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/gofiber/swagger"

	_ "github.com/Xapsiel/bpla_dashboard/docs"
	"github.com/Xapsiel/bpla_dashboard/internal/config"
	httpv1 "github.com/Xapsiel/bpla_dashboard/internal/entrypoint"
	"github.com/Xapsiel/bpla_dashboard/internal/repository"
	"github.com/Xapsiel/bpla_dashboard/internal/service"
)

func InitLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return logger
}

// @title BPLA Dashboard API
// @version 1.0
// @description API для дашборда БПЛА
// @schemes http https
// @host 127.0.0.1:8080
// @BasePath /
func main() {
	logger := InitLogger()
	slog.SetDefault(logger)
	configPath := flag.String("c", "config/config.yaml", "The path to the configuration file")
	flag.Parse()

	cfg, err := config.New(*configPath)
	if err != nil {
		slog.Error("unable to read config:", err.Error())
		os.Exit(1)
	}
	db, err := repository.NewPostgresDB(cfg.DatabaseConfig)
	if err != nil {
		slog.Error("unable to connect to database:", err.Error())
		os.Exit(1)
	}
	defer db.Close()
	repo := repository.NewRepository(db)

	app := fiber.New(fiber.Config{
		BodyLimit: cfg.MaxUploadSize * 1024 * 1024,
	})

	// Swagger UI endpoint
	app.Get("/swagger/*", fiberSwagger.New())
	service := service.New(repo, cfg.OidcConfig)
	router := httpv1.New(httpv1.Config{
		Repo:         repo,
		Domain:       cfg.Domain,
		Service:      &service,
		IsProduction: cfg.IsProduction,
	})
	router.Routes(app)
	log.Fatal(app.Listen(":" + cfg.HostConfig.Port))
}
