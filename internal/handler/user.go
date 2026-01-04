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
// @Tags 			Users
// @Accept 			json
// @Produce 		json
// @Param 			user body dto.CreateUserRequest true "User info"
// @Success 		201 {object} models.User
// @Failure 		400 {object} map[string]string
// @Router 			/users [post]
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

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

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

func NoAPI(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "route not found"})
}

func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "OK"})
}
