package httpv1

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func (r *Router) DistrictGeoJSONHandler(ctx *fiber.Ctx) error {
	name := ctx.Query("name")
	switch name {
	case "":
		res, err := r.repo.GetAllDistrictsGeoJSONHandler(context.Background())
		if err != nil {
			slog.Error("failed to get districts geojson", "error", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(r.NewErrorResponse(fiber.StatusInternalServerError, "Ошибка при получении данных по районам"))
		}
		return ctx.Status(fiber.StatusOK).JSON(r.NewSuccessResponse(res, ""))
	default:
		res, err := r.repo.GetDistrictGeoJSON(context.Background(), name)
		if err != nil {
			slog.Error("failed to get district geojson", "name", name, "error", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(r.NewErrorResponse(fiber.StatusInternalServerError, "Ошибка при получении данных по району"))
		}
		return ctx.Status(fiber.StatusOK).JSON(r.NewSuccessResponse(res, ""))

	}
}
