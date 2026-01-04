package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/sudo-hassan-zahid/go-api-server/internal/auth"
	dto "github.com/sudo-hassan-zahid/go-api-server/internal/dto"
	appErrors "github.com/sudo-hassan-zahid/go-api-server/internal/errors"
	"github.com/sudo-hassan-zahid/go-api-server/internal/logger"
	"github.com/sudo-hassan-zahid/go-api-server/internal/service"
	"github.com/sudo-hassan-zahid/go-api-server/utils"
	"gorm.io/gorm"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
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
func (h *AuthHandler) CreateUser(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to parse request body")
		return appErrors.HandleError(c, appErrors.ErrBadRequest)
	}

	if ok := utils.ValidateStruct(c, &req); !ok {
		logger.Log.Warn().Msg("Validation failed")
		return nil
	}

	user, err := h.service.CreateUser(req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to create user")
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return appErrors.HandleError(c, appErrors.ErrEmailAlreadyExists)
		}
		return appErrors.HandleError(c, err)
	}
	logger.Log.Info().Msg("User created successfully")
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
func (h *AuthHandler) LoginUser(c *fiber.Ctx) error {
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
