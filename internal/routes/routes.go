package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nattakan-n/ai-training-backend/internal/handlers"
	"github.com/nattakan-n/ai-training-backend/internal/middleware"
)

func Register(app *fiber.App, authHandler *handlers.AuthHandler, jwtSecret string) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	auth := v1.Group("")
	auth.Post("/login", authHandler.Login)

	protected := v1.Group("")
	protected.Use(middleware.JWTAuth(jwtSecret))
	protected.Get("/me", authHandler.Me)
	protected.Put("/me", authHandler.UpdateProfile)
}


