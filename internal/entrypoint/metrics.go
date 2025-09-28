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
		slog.Error("failed to get metrics", "reg_id", regID, "year", year, "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(r.NewErrorResponse(fiber.StatusInternalServerError, "Ошибка при получении метрик"))
	}
	return ctx.Status(fiber.StatusOK).JSON(r.NewSuccessResponse(metrics, ""))
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
			slog.Error("failed to get metrics for region", "region_id", *regID.Gid, "year", year, "error", err)
			continue
		}
		metrics = append(metrics, &m)
	}
	return ctx.Status(fiber.StatusOK).JSON(r.NewSuccessResponse(metrics, ""))

}
