package httpv1

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (r *Router) GenerateAuthURLHandler(ctx *fiber.Ctx) error {

	return ctx.JSON(fiber.Map{
		"res": r.service.UserService.GetAuthURL("a"),
	})
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
	r.service.UserService.ExchangeCode(code)
	return ctx.JSON(fiber.Map{})
}
