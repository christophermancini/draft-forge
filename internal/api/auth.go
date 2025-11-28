package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/yourusername/draft-forge/internal/auth"
	"github.com/yourusername/draft-forge/internal/models"
)

type authHandler struct {
	service      *auth.Service
	tokenManager *auth.TokenManager
	userStore    auth.Store
}

func NewAuthHandler(service *auth.Service, tokenManager *auth.TokenManager, userStore auth.Store) *authHandler {
	return &authHandler{
		service:      service,
		tokenManager: tokenManager,
		userStore:    userStore,
	}
}

func (h *authHandler) Register(app fiber.Router) {
	app.Get("/auth/github/start", h.start)
	app.Get("/auth/github/callback", h.callback)

	protected := app.Group("", AuthMiddleware(h.tokenManager, h.userStore))
	protected.Get("/me", h.me)
}

func (h *authHandler) start(c *fiber.Ctx) error {
	start, err := h.service.StartAuth()
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"data": start,
	})
}

func (h *authHandler) callback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	result, err := h.service.CompleteAuth(c.Context(), code, state)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrMissingCode), errors.Is(err, auth.ErrInvalidState):
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		default:
			return err
		}
	}

	return c.JSON(fiber.Map{"data": result})
}

func (h *authHandler) me(c *fiber.Ctx) error {
	userIDVal := c.Locals("user_id")
	userID, ok := userIDVal.(int64)
	if !ok || userID <= 0 {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	user, err := h.userStore.GetUserByID(c.Context(), userID)
	if err != nil {
		return err
	}
	c.Locals("github_token", user.AccessToken)
	return c.JSON(fiber.Map{"data": models.User{
		ID:        user.ID,
		GitHubID:  user.GitHubID,
		Username:  user.Username,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	}})
}
