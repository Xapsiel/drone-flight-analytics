package main

import (
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func InitLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return logger
}

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
	engine := html.New("./web/views", ".gohtml")
	engine.AddFunc(
		"scoreFilled", func(score int) []struct{} {
			return make([]struct{}, score)
		},
	)
	engine.AddFunc(
		"scoreEmpty", func(score int) []struct{} {
			return make([]struct{}, 5-score)
		},
	)

	engine.Reload(!cfg.IsProduction)

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
		// TODO: custom ErrorHandler:
	})

	router := httpv1.New(httpv1.Config{
		Repo:       repo,
		Domain:     cfg.Domain,
		TileServer: cfg.TileServer,
	})
	router.Routes(app)

	log.Fatal(app.Listen(":" + cfg.HostConfig.Port))
}
