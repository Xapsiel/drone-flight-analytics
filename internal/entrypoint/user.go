package httpv1

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

// GenerateAuthURLHandler
// @Summary Сгенерировать ссылку авторизации
// @Description Возвращает URL для авторизации через провайдер OIDC
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} httpv1.APIResponse
// @Router /user/gen_auth_url [get]
func (r *Router) GenerateAuthURLHandler(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(r.NewSuccessResponse(fiber.Map{
		"auth_url": r.service.UserService.GetAuthURL("a"),
	}, ""))
}

// RedirectAuthURLHandler
// @Summary Обработать редирект от OIDC
// @Description Обменивает код авторизации на токены
// @Tags user
// @Accept json
// @Produce json
// @Param state query string false "State"
// @Param session_state query string false "Session state"
// @Param iss query string false "Issuer"
// @Param code query string true "Authorization code"
// @Success 200 {object} httpv1.APIResponse
// @Failure 500 {object} httpv1.APIResponse
// @Router /user/redirect [get]
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
