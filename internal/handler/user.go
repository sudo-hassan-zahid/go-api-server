package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	appErrors "github.com/sudo-hassan-zahid/go-api-server/internal/errors"
	"github.com/sudo-hassan-zahid/go-api-server/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// GetAllUsers 	 godoc
// @Summary      Get all users
// @Description  Returns a list of all existing users
// @Tags         Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200 {array} models.User "List of users"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

// GetUserByID 	 godoc
// @Summary      Get user by ID
// @Description  Returns a single user by their ID
// @Tags         Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "User UUID"
// @Success      200 {object} models.User "User found"
// @Failure      400 {object} map[string]string "Invalid ID"
// @Failure      404 {object} map[string]string "User not found"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return appErrors.HandleError(c, appErrors.ErrBadRequest)
	}

	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		return appErrors.HandleError(c, err)
	}
	return c.JSON(user)
}
