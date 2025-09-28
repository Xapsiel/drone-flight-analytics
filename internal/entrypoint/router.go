package httpv1

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"

	"github.com/Xapsiel/bpla_dashboard/internal/model"
	"github.com/Xapsiel/bpla_dashboard/internal/service"
)

type Repository interface {
	GetDistrictGeoJSON(ctx context.Context, name string) (*model.DistrictGeoJSON, error)
	GetAllDistrictsGeoJSONHandler(ctx context.Context) ([]model.DistrictGeoJSON, error)
	SaveFileInfo(background context.Context, mf model.File, valid_count int, error_count int) (int, error)
	GetMetrics(ctx context.Context, id int, year int) (model.Metrics, error)
	GetRegions(ctx context.Context) []model.District
	GetDistrictsMVT(ctx context.Context, z, x, y int) ([]byte, error)
}
type Router struct {
	repo         Repository
	domain       string
	isProduction bool
	service      *service.Service
}

type Config struct {
	Repo         Repository
	Service      *service.Service
	Domain       string
	IsProduction bool
}

func New(cfg Config) *Router {
	return &Router{
		repo:         cfg.Repo,
		domain:       cfg.Domain,
		isProduction: cfg.IsProduction,
		service:      cfg.Service,
	}
}

func (r *Router) Routes(app fiber.Router) {
	// Настройка CORS для фронтенда
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173,http://127.0.0.1:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	app.Static("assets", "web/assets")
	app.Static("", "web/assets")
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	app.Get("/dashboard", monitor.New())
	app.Get("/tiles/:z/:x/:y.mvt", r.GetTileMVT)
	district := app.Group("/district")
	if r.isProduction {
		district.Use(r.RoleMiddleware("admin", "analytics"))
	}
	district.Get("/", r.DistrictGeoJSONHandler)

	user := app.Group("/user")
	user.Get("/gen_auth_url", r.GenerateAuthURLHandler)
	user.Get("/redirect", r.RedirectAuthURLHandler)
	user.Get("/me", r.GetCurrentUserHandler)
	user.Post("/logout", r.LogoutHandler)
	user.Post("/refresh", r.RefreshTokenHandler)

	// Эндпоинт для обработки callback от фронтенда
	app.Get("/auth/callback", r.AuthCallbackHandler)

	crawler := app.Group("/crawler")
	crawler.Post("/upload", r.UploadFileHandler)

	metrics := app.Group("/metrics")
	metrics.Get("/", r.GetMetrics)
	metrics.Get("/all", r.GetAllMetrics)

}

func (r *Router) NewPage() *model.Page {
	return &model.Page{
		Domain: r.domain,
		Year:   time.Now().Year(),
	}
}

func (r *Router) NewErrorPage(err error) *model.Page {
	return &model.Page{
		Error:  err.Error(),
		Domain: r.domain,
		Year:   time.Now().Year(),
	}
}

func (r *Router) GetTopByHandler(ctx *fiber.Ctx) error {
	filter := ctx.Query("criteria", "none")
	if filter == "" {
		filter = "flight_frequency"
	}
	switch filter {
	case "flight_frequency":
	case "avg_flight_time":
	case "flight_count":
	case "flight_duration":
	}
	return ctx.JSON(fiber.Map{})
}
