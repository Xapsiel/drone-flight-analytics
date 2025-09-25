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

// http://localhosiss=http%3A%2F%2F51.250.101.118%3A8080%2Frealms%2Fdashboard&code=602e157f-d285-4d73-ae45-07ecf07600ec.057a9b02-963d-455f-93f6-d1153a838ab7.42c377f8-540f-43be-a459-bc492408ede6
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
