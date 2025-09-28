package httpv1

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (r *Router) RoleMiddleware(allowedRoles ...string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: missing token",
			})
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		user, err := r.service.UserService.GetCurrentUser(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: " + err.Error(),
			})
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
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden: insufficient role",
			})
		}

		return c.Next()
	}
}
