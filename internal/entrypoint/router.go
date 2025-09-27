package httpv1

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"

	"github.com/Xapsiel/bpla_dashboard/internal/model"
	"github.com/Xapsiel/bpla_dashboard/internal/service"
)

type Repository interface {
	GetDistrictGeoJSON(ctx context.Context, name string) (*model.DistrictGeoJSON, error)
	GetAllDistrictsGeoJSONHandler(ctx context.Context) ([]model.DistrictGeoJSON, error)
	SaveFileInfo(background context.Context, mf model.File, valid_count int, error_count int) (int, error)
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
	app.Static("", "web/assets")
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	app.Get("/dashboard", monitor.New())
	district := app.Group("/district")
	if r.isProduction {
		district.Use(r.RoleMiddleware("admin", "analytics"))
	}
	district.Get("/", r.DistrictGeoJSONHandler)
	district.Get("/top", r.GetTopByHandler)
	user := app.Group("/user")
	user.Get("/gen_auth_url", r.GenerateAuthURLHandler)
	user.Get("/redirect", r.RedirectAuthURLHandler)

	crawler := app.Group("/crawler")
	crawler.Post("/upload", r.UploadFileHandler)

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
