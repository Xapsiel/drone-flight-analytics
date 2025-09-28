package httpv1

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func (r *Router) GenerateAuthURLHandler(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(r.NewSuccessResponse(fiber.Map{
		"auth_url": r.service.UserService.GetAuthURL("a"),
	}, ""))
}

func (r *Router) RedirectAuthURLHandler(ctx *fiber.Ctx) error {
	state := ctx.Query("state")
	session_state := ctx.Query("session_state")
	iss := ctx.Query("iss")
	code := ctx.Query("code")
	fmt.Println("state:", state)
	fmt.Println("session_state:", session_state)
	fmt.Println("iss:", iss)
	fmt.Println("code:", code)
	_, err := r.service.UserService.ExchangeCode(code)
	if err != nil {
		slog.Error("failed to exchange auth code", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(r.NewErrorResponse(fiber.StatusInternalServerError, "Ошибка авторизации"))
	}
	return ctx.Status(fiber.StatusOK).JSON(r.NewSuccessResponse(fiber.Map{}, ""))
}
