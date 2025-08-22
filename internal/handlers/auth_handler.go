package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nattakan-n/ai-training-backend/internal/usecases"
)

type AuthHandler struct {
	service *usecases.AuthService
}

func NewAuthHandler(service *usecases.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type updateProfileRequest struct {
	Name string `json:"name"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	token, user, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		log.Printf("error Login: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}
	return c.JSON(fiber.Map{"token": token, "email": user.Email, "name": user.Name})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID, _ := c.Locals("userID").(string)
	user, err := h.service.GetProfile(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(fiber.Map{"id": user.ID, "email": user.Email, "name": user.Name})
}

func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, _ := c.Locals("userID").(string)
	var req updateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	user, err := h.service.UpdateProfile(userID, req.Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"id": user.ID, "email": user.Email, "name": user.Name})
}


