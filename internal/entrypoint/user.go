package httpv1

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (r *Router) GenerateAuthURLHandler(ctx *fiber.Ctx) error {

	return ctx.JSON(fiber.Map{
		"res": r.service.UserService.GetAuthURL("a"),
	})
}

func (r *Router) RedirectAuthURLHandler(ctx *fiber.Ctx) error {
	// Получаем параметры от Keycloak (некоторые могут быть не использованы)
	_ = ctx.Query("state")
	_ = ctx.Query("session_state")
	_ = ctx.Query("iss")
	code := ctx.Query("code")

	// Валидация обязательных параметров
	if code == "" {
		slog.Error("Authorization code is missing")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Authorization code is required",
		})
	}

	// Обмен кода на токен и получение информации о пользователе
	user, err := r.service.UserService.ExchangeCode(code)
	if err != nil {
		slog.Error("Failed to exchange code", "error", err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to authenticate user",
		})
	}

	// Сохранение информации о пользователе в сессии или базе данных
	// Здесь можно добавить логику сохранения пользователя в БД

	// Создание JWT токена для сессии (опционально)
	sessionToken, err := r.service.UserService.CreateSessionToken(user)
	if err != nil {
		slog.Error("Failed to create session token", "error", err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create session",
		})
	}

	// Установка cookies или редирект на фронтенд с токеном
	ctx.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    sessionToken,
		HTTPOnly: true,
		Secure:   r.isProduction,
		SameSite: "Lax",
		MaxAge:   3600, // 1 час
	})

	// Редирект на фронтенд с токеном в URL (если нужно)
	frontendURL := fmt.Sprintf("%s/auth/callback?token=%s", r.domain, sessionToken)

	return ctx.Redirect(frontendURL, fiber.StatusFound)
}

// GetCurrentUserHandler возвращает информацию о текущем пользователе
func (r *Router) GetCurrentUserHandler(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	slog.Info("GetCurrentUserHandler called", "authHeader", authHeader)

	if authHeader == "" {
		slog.Error("No authorization header")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: missing token",
		})
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	slog.Info("Token extracted", "tokenLength", len(tokenStr))

	user, err := r.service.UserService.GetCurrentUser(tokenStr)
	if err != nil {
		slog.Error("Failed to get current user", "error", err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: " + err.Error(),
		})
	}

	slog.Info("User retrieved successfully", "username", user.Username, "roles", user.Roles)
	return ctx.JSON(fiber.Map{
		"user": user,
	})
}

// LogoutHandler обрабатывает выход пользователя
func (r *Router) LogoutHandler(ctx *fiber.Ctx) error {
	// Очистка cookie
	ctx.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   r.isProduction,
		SameSite: "Lax",
		MaxAge:   -1, // Удалить cookie
	})

	return ctx.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

// RefreshTokenHandler обновляет токен пользователя
func (r *Router) RefreshTokenHandler(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: missing token",
		})
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	user, err := r.service.UserService.GetCurrentUser(tokenStr)
	if err != nil {
		slog.Error("Failed to validate token for refresh", "error", err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: " + err.Error(),
		})
	}

	// Создание нового токена
	newToken, err := r.service.UserService.CreateSessionToken(*user)
	if err != nil {
		slog.Error("Failed to create new session token", "error", err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to refresh token",
		})
	}

	// Установка нового cookie
	ctx.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    newToken,
		HTTPOnly: true,
		Secure:   r.isProduction,
		SameSite: "Lax",
		MaxAge:   3600, // 1 час
	})

	return ctx.JSON(fiber.Map{
		"token": newToken,
		"user":  user,
	})
}

// AuthCallbackHandler обрабатывает callback от фронтенда
func (r *Router) AuthCallbackHandler(ctx *fiber.Ctx) error {
	token := ctx.Query("token")
	if token == "" {
		// Если токена нет, редиректим на главную страницу фронтенда
		return ctx.Redirect("http://localhost:5173", fiber.StatusFound)
	}

	// Устанавливаем токен в cookie и редиректим на фронтенд
	ctx.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   r.isProduction,
		SameSite: "Lax",
		MaxAge:   3600, // 1 час
	})

	// Редирект на фронтенд
	return ctx.Redirect(fmt.Sprintf("http://localhost:5173?token=%s", token), fiber.StatusFound)
}
