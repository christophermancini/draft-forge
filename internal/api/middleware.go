package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/yourusername/draft-forge/internal/auth"
)

func AuthMiddleware(tokenMgr *auth.TokenManager, userStore auth.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := tokenMgr.ParseAccessToken(token)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}

		c.Locals("user_id", claims.UserID)

		if userStore != nil {
			if user, err := userStore.GetUserByID(c.Context(), claims.UserID); err == nil {
				c.Locals("github_token", user.AccessToken)
			}
		}

		return c.Next()
	}
}
