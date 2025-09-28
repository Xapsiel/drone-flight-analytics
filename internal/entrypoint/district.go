package httpv1

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (r *Router) DistrictGeoJSONHandler(ctx *fiber.Ctx) error {
	name := ctx.Query("name")
	switch name {
	case "":
		res, err := r.repo.GetAllDistrictsGeoJSONHandler(context.Background())
		if err != nil {
			return ctx.JSON(fiber.Map{
				"code": fiber.StatusInternalServerError,
				"msg":  err.Error(),
			})
		}
		return ctx.JSON(fiber.Map{
			"code": fiber.StatusOK,
			"data": res,
		})
	default:
		res, err := r.repo.GetDistrictGeoJSON(context.Background(), name)
		if err != nil {
			return ctx.JSON(fiber.Map{
				"code": fiber.StatusInternalServerError,
				"msg":  err.Error(),
			})
		}
		return ctx.JSON(fiber.Map{
			"code": fiber.StatusOK,
			"data": res,
		})

	}
	return ctx.JSON(fiber.Map{
		"code": fiber.StatusOK,
	})
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
		return ctx.Status(fiber.StatusBadRequest).JSON(NewErrorResponse(fiber.StatusBadRequest, "Некорректный параметр z"))
	}
	x, err := strconv.Atoi(ctx.Params("x"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(NewErrorResponse(fiber.StatusBadRequest, "Некорректный параметр x"))
	}
	y, err := strconv.Atoi(ctx.Params("y"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(NewErrorResponse(fiber.StatusBadRequest, "Некорректный параметр y"))
	}
	tile, err := r.repo.GetDistrictsMVT(context.Background(), z, x, y)
	if err != nil {
		slog.Error("mvt generation failed", "z", z, "x", x, "y", y, "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(NewErrorResponse(fiber.StatusInternalServerError, "Ошибка генерации тайла"))
	}
	if len(tile) == 0 {
		return ctx.SendStatus(fiber.StatusNoContent)
	}
	ctx.Set("Content-Type", "application/x-protobuf")
	return ctx.Send(tile)
}
