package httpv1

import (
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (r *Router) RoleMiddleware(allowedRoles ...string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(r.NewErrorResponse(fiber.StatusUnauthorized, "Unauthorized: missing token"))
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		user, err := r.service.UserService.GetUserInfo("", tokenStr)
		if err != nil {
			slog.Error("failed to parse user token", "error", err)
			return c.Status(fiber.StatusUnauthorized).JSON(r.NewErrorResponse(fiber.StatusUnauthorized, "Unauthorized: invalid token"))
		}

		hasRole := false
		for _, role := range user.Roles {
			for _, allowed := range allowedRoles {
				if role == allowed {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(r.NewErrorResponse(fiber.StatusForbidden, "Forbidden: insufficient role"))
		}

		return c.Next()
	}
}
