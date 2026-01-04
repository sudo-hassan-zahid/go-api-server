package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	dto "github.com/sudo-hassan-zahid/go-api-server/internal/dto"
	"github.com/sudo-hassan-zahid/go-api-server/internal/service"
	"github.com/sudo-hassan-zahid/go-api-server/utils"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// CreateUser 		godoc
// @Summary 		Create a new user
// @Description 	Creates a user with email and password
// @Tags 			Auth
// @Accept 			json
// @Produce 		json
// @Param 			user body dto.CreateUserRequest true "User info"
// @Success 		201 {object} models.User
// @Failure 		400 {object} map[string]string
// @Router 			/auth/signup [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if ok := utils.ValidateStruct(c, &req); !ok {
		return nil
	}

	user, err := h.service.CreateUser(req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// LoginUser 		godoc
// @Summary 		Login an existing user
// @Description 	Logins an existing user with email and password
// @Tags 			Auth
// @Accept 			json
// @Produce 		json
// @Param 			user body dto.LoginUserRequest true "User info"
// @Success 		200 {object} map[string]interface{}
// @Failure 		401 {object} map[string]string
// @Router 			/auth/login [post]
func (h *UserHandler) LoginUser(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.service.LoginUser(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "login successful", "user": user})
}

// GetAllUsers 		godoc
// @Summary 		Get all users
// @Description 	Get all existing users
// @Tags 			Users
// @Accept 			json
// @Produce 		json
// @Success 		200 {array} []models.User
// @Failure 		500 {object} map[string]string
// @Router 			/users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

// @Summary 		Get a user by ID
// @Description 	Get a user by their ID
// @Tags 			Users
// @Accept 			json
// @Produce 		json
// @Param 			id path uint true "User ID"
// @Success 		200 {object} models.User
// @Failure 		400 {object} map[string]string
// @Failure 		404 {object} map[string]string
// @Router 			/users/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}
	return c.JSON(user)
}
