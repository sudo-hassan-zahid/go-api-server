package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sudo-hassan-zahid/go-api-server/internal/auth"
	"github.com/sudo-hassan-zahid/go-api-server/internal/dto"
	"github.com/sudo-hassan-zahid/go-api-server/internal/service"
	"github.com/sudo-hassan-zahid/go-api-server/utils"
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
		return err
	}

	if ok := utils.ValidateStruct(c, &req); !ok {
		return nil
	}

	user, err := h.service.CreateUser(req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		return err
	}

	accessToken, err := auth.GenerateAccessToken(user.ID.String(), user.Role)
	if err != nil {
		return err
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID.String(), user.Role)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(dto.SignupResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.ID.String(),
		Role:         user.Role,
	})
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
		return err
	}

	if ok := utils.ValidateStruct(c, &req); !ok {
		return nil
	}

	user, err := h.service.LoginUser(req.Email, req.Password)
	if err != nil {
		return err
	}

	accessToken, err := auth.GenerateAccessToken(user.ID.String(), user.Role)
	if err != nil {
		return err
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID.String(), user.Role)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.LoginUserResponse{
		UserID:       user.ID.String(),
		UserRole:     user.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
