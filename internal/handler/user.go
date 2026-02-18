package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sudo-hassan-zahid/go-api-server/internal/dto"
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
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Success      200 {array} models.User "List of users"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return err
	}
	return c.JSON(users)
}

// GetUserByID 	 godoc
// @Summary      Get user by ID
// @Description  Returns a single user by their ID
// @Tags         Users
// @Security     Bearer
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
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid user id format",
		})
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

// UpdateUser 	 godoc
// @Summary      Update a user
// @Description  Updates a user with given ID
// @Tags         Users
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        id path string true "User UUID"
// @Param        user body dto.UpdateUserRequest true "User info"
// @Success      200 {object} dto.SuccessResponse "User updated successfully"
// @Failure      400 {object} map[string]string "Bad request / validation error"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /users/{id} [patch]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	userID := strings.TrimSpace(c.Params("id"))
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "user id is required",
		})
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid user id format",
		})
	}

	_, err = h.service.GetUserByID(id)
	if err != nil {
		return err
	}

	var updateUser dto.UpdateUserRequest
	if err := c.BodyParser(&updateUser); err != nil {
		return err
	}

	err = h.service.UpdateUser(id, updateUser)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

// DeleteUser 	 godoc
// @Summary      Delete a user
// @Description  Deletes a user with given ID
// @Tags         Users
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        id path string true "User UUID"
// @Success      200 {object} map[string]string "User deleted successfully"
// @Failure      400 {object} map[string]string "Bad request / validation error"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	panic("not implemented")
}
