package handlers

import (
	"log"
	"module_6/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// SignUp godoc
// @Summary Register new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{username=string,password=string} true "User credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/sign-up [post]
func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] SignUp: invalid request body: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	log.Printf("[INFO] SignUp attempt: username=%s\n", req.Username)
	if err := h.service.SignUp(req.Username, req.Password); err != nil {
		log.Printf("[ERROR] SignUp failed for username=%s: %v\n", req.Username, err)
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("[INFO] User created: username=%s\n", req.Username)
	return c.JSON(fiber.Map{"message": "user created"})
}

// SignIn godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{username=string,password=string} true "User credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/sign-in [post]
func (h *AuthHandler) SignIn(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] SignIn: invalid request body: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	log.Printf("[INFO] SignIn attempt: username=%s\n", req.Username)
	token, err := h.service.SignIn(req.Username, req.Password)
	if err != nil {
		log.Printf("[ERROR] SignIn failed for username=%s: %v\n", req.Username, err)
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("[INFO] User signed in: username=%s\n", req.Username)
	return c.JSON(fiber.Map{"token": token})
}
