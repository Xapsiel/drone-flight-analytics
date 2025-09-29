package httpv1

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

func (r *Router) GetDistricsIDs(ctx *fiber.Ctx) error {
	districts := r.repo.GetRegions(context.Background())
	if len(districts) == 0 {
		return ctx.Status(fiber.StatusInternalServerError).JSON(r.NewErrorResponse(fiber.StatusInternalServerError, "Ошибка при получении регионов"))
	}
	return ctx.Status(fiber.StatusOK).JSON(r.NewSuccessResponse(districts, ""))

}
