package httpv1

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/Xapsiel/bpla_dashboard/internal/model"
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
func (r *Router) GetAllMetrics(ctx *fiber.Ctx) error {
	year, err := strconv.Atoi(ctx.Query("year"))
	if err != nil {
		year = 2025
	}
	reg := r.repo.GetRegions(context.Background())
	metrics := []*model.Metrics{}
	for _, regID := range reg {
		m, err := r.repo.GetMetrics(context.Background(), *regID.Gid, year)
		if err != nil {
			slog.Info("error with getting metrics: %v", err)
			continue
		}
		metrics = append(metrics, &m)
	}
	return ctx.JSON(fiber.Map{
		"metrics": metrics,
	})

}
