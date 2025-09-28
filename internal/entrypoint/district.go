package httpv1

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// DistrictGeoJSONHandler
// @Summary Получить GeoJSON районов
// @Description Возвращает GeoJSON по всем районам или конкретному району при переданном параметре id
// @Tags district
// @Accept json
// @Produce json
// @Param id query int false "ID региона"
// @Success 200 {object} httpv1.APIResponse
// @Failure 500 {object} httpv1.APIResponse
// @Router /district [get]
func (r *Router) DistrictGeoJSONHandler(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		id = 0
	}
	switch id {
	case 0:
		res, err := r.repo.GetAllDistrictsGeoJSONHandler(context.Background())
		if err != nil {
			slog.Error("failed to get districts geojson", "error", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(r.NewErrorResponse(fiber.StatusInternalServerError, "Ошибка при получении данных по районам"))
		}
		return ctx.Status(fiber.StatusOK).JSON(r.NewSuccessResponse(res, ""))
	default:
		res, err := r.repo.GetDistrictGeoJSON(context.Background(), id)
		if err != nil {
			slog.Error("failed to get district geojson", "id", id, "error", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(r.NewErrorResponse(fiber.StatusInternalServerError, "Ошибка при получении данных по району"))
		}
		return ctx.Status(fiber.StatusOK).JSON(r.NewSuccessResponse(res, ""))

	}
}
