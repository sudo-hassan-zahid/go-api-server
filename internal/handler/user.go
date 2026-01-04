package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sudo-hassan-zahid/go-api-server/internal/auth"
	dto "github.com/sudo-hassan-zahid/go-api-server/internal/dto"
	appErrors "github.com/sudo-hassan-zahid/go-api-server/internal/errors"
	"github.com/sudo-hassan-zahid/go-api-server/internal/service"
	"github.com/sudo-hassan-zahid/go-api-server/utils"
	"gorm.io/gorm"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// CreateUser 	 godoc
// @Summary      Create a new user
// @Description  Creates a new user with email, password, first name, and last name
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user body dto.CreateUserRequest true "User info"
// @Success      201 {object} models.User "Created user"
// @Failure      400 {object} map[string]string "Bad request / validation error"
// @Failure      409 {object} map[string]string "Email already exists"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /auth/signup [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return appErrors.HandleError(c, appErrors.ErrBadRequest)
	}

	if ok := utils.ValidateStruct(c, &req); !ok {
		return nil
	}

	user, err := h.service.CreateUser(req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return appErrors.HandleError(c, appErrors.ErrEmailAlreadyExists)
		}
		return appErrors.HandleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// LoginUser 	 godoc
// @Summary      Login an existing user
// @Description  Logins an existing user using email and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user body dto.LoginUserRequest true "User credentials"
// @Success      200 {object} dto.LoginUserResponse "Login successful, returns user object"
// @Failure      400 {object} map[string]string "Bad request / validation error"
// @Failure      401 {object} map[string]string "Invalid credentials"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /auth/login [post]
func (h *UserHandler) LoginUser(c *fiber.Ctx) error {
	var req dto.LoginUserRequest
	if err := c.BodyParser(&req); err != nil {
		return appErrors.HandleError(c, appErrors.ErrBadRequest)
	}

	if ok := utils.ValidateStruct(c, &req); !ok {
		return nil
	}

	user, err := h.service.LoginUser(req.Email, req.Password)
	if err != nil {
		return appErrors.HandleError(c, err)
	}

	accessToken, err := auth.GenerateAccessToken(user.ID.String(), user.Role)
	if err != nil {
		return appErrors.HandleError(c, err)
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return appErrors.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(dto.LoginUserResponse{
		UserID:       user.ID.String(),
		UserRole:     user.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
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
