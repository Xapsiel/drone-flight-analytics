package httpv1

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	fiberSwagger "github.com/gofiber/swagger"

	"github.com/Xapsiel/bpla_dashboard/internal/model"
	"github.com/Xapsiel/bpla_dashboard/internal/service"
)

type Repository interface {
	GetDistrictGeoJSON(ctx context.Context, id int) ([]byte, error)
	GetAllDistrictsGeoJSONHandler(ctx context.Context) ([]byte, error)
	GetDistrictsMVT(ctx context.Context, z, x, y int) ([]byte, error)
	SaveFileInfo(background context.Context, mf model.File, valid_count int, error_count int) (int, error)
	GetMetrics(ctx context.Context, id int, year int) (model.Metrics, error)
	GetRegions(ctx context.Context) []model.District
	GetFile(ctx context.Context, id int) (model.File, error)
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
	app.Static("assets", "web/assets")
	// Раздача статических тестовых данных (.data)
	app.Static(".data", ".data")
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	app.Get("/swagger/*", fiberSwagger.New())

	app.Get("/dashboard", monitor.New())
	//district := app.Group("/district")
	//district.Get("/", r.DistrictGeoJSONHandler)
	user := app.Group("/user")
	user.Use()
	user.Get("/gen_auth_url", r.GenerateAuthURLHandler)
	user.Get("/redirect", r.RedirectAuthURLHandler)

	crawler := app.Group("/crawler")
	crawler.Use(r.RoleMiddleware("admin"))
	crawler.Post("/upload", r.UploadFileHandler)
	crawler.Get("/status", r.CheckFileStatus)

	metrics := app.Group("/metrics")
	metrics.Use(r.RoleMiddleware("admin", "analytics"))
	metrics.Get("/", r.GetMetrics)
	metrics.Get("/all", r.GetAllMetrics)

	app.Get("/tiles/:z/:x/:y.mvt", r.GetTileMVT)
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
