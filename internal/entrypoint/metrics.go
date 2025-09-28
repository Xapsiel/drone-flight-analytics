package httpv1

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (r *Router) GetMetrics(ctx *fiber.Ctx) error {

	regID, err := strconv.Atoi(ctx.Query("reg_id"))
	if err != nil {
		regID = 0
	}
	year, err := strconv.Atoi(ctx.Query("year"))
	if err != nil {
		year = 2025
	}
	metrics, err := r.repo.GetMetrics(context.Background(), regID, year)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(fiber.Map{
		"metrics": metrics,
	})
}
