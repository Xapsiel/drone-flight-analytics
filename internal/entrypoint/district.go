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
	if id == 0 {
		res, err := r.repo.GetAllDistrictsGeoJSONHandler(context.Background())
		if err != nil {
			slog.Error("failed to get districts geojson", "error", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(r.NewErrorResponse(fiber.StatusInternalServerError, "Ошибка при получении данных по районам"))
		}
		ctx.Type("json")
		return ctx.Status(fiber.StatusOK).Send(res)
	}
	res, err := r.repo.GetDistrictGeoJSON(context.Background(), id)
	if err != nil {
		slog.Error("failed to get district geojson", "id", id, "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(r.NewErrorResponse(fiber.StatusInternalServerError, "Ошибка при получении данных по району"))
	}
	ctx.Type("json")
	return ctx.Status(fiber.StatusOK).Send(res)
}

// GetTileMVT
// @Summary Векторный тайл MVT по z/x/y
// @Description Возвращает тайл районов в формате Mapbox Vector Tile (protobuf)
// @Tags district
// @Produce application/x-protobuf
// @Param z path int true "Zoom"
// @Param x path int true "Tile X"
// @Param y path int true "Tile Y"
// @Success 200 {string} string "binary mvt"
// @Failure 500 {object} httpv1.APIResponse
// @Router /tiles/{z}/{x}/{y}.mvt [get]
func (r *Router) GetTileMVT(ctx *fiber.Ctx) error {
	z, err := strconv.Atoi(ctx.Params("z"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(r.NewErrorResponse(fiber.StatusBadRequest, "Некорректный параметр z"))
	}
	x, err := strconv.Atoi(ctx.Params("x"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(r.NewErrorResponse(fiber.StatusBadRequest, "Некорректный параметр x"))
	}
	y, err := strconv.Atoi(ctx.Params("y"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(r.NewErrorResponse(fiber.StatusBadRequest, "Некорректный параметр y"))
	}
	tile, err := r.repo.GetDistrictsMVT(context.Background(), z, x, y)
	if err != nil {
		slog.Error("mvt generation failed", "z", z, "x", x, "y", y, "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(r.NewErrorResponse(fiber.StatusInternalServerError, "Ошибка генерации тайла"))
	}
	if len(tile) == 0 {
		return ctx.SendStatus(fiber.StatusNoContent)
	}
	ctx.Set("Content-Type", "application/x-protobuf")
	return ctx.Send(tile)
}
