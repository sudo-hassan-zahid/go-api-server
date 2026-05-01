package handler

import (
	"github.com/gofiber/fiber/v2"
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
// @Success      429 {object} map[string]string "Too many requests"
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
// @Failure      429 {object} map[string]string "Too many requests"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := uuidParam(c, "id")
	if err != nil {
		return err
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
// @Success      200 {object} map[string]string "User updated successfully"
// @Failure      400 {object} map[string]string "Bad request / validation error"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      429 {object} map[string]string "Too many requests"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /users/{id} [patch]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := uuidParam(c, "id")
	if err != nil {
		return err
	}

	_, err = h.service.GetUserByID(id)
	if err != nil {
		return err
	}

	updateUser, ok, err := bindAndValidate[dto.UpdateUserRequest](c)
	if err != nil {
		return err
	}
	if !ok {
		return nil
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
// @Failure      404 {object} map[string]string "User not found"
// @Failure      429 {object} map[string]string "Too many requests"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := uuidParam(c, "id")
	if err != nil {
		return err
	}

	_, err = h.service.GetUserByID(id)
	if err != nil {
		return err
	}

	err = h.service.DeleteUser(id)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}
