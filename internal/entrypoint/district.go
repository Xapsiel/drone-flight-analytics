package httpv1

import (
	"context"

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
