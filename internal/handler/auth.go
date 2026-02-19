package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sudo-hassan-zahid/go-api-server/internal/auth"
	"github.com/sudo-hassan-zahid/go-api-server/internal/dto"
	"github.com/sudo-hassan-zahid/go-api-server/utils"
)

type AuthHandler struct {
	service *auth.Service
}

func NewAuthHandler(s *auth.Service) *AuthHandler {
	return &AuthHandler{service: s}
}

// CreateUser 	 godoc
// @Summary      Create a new user
// @Description  Creates a new user with email, password, first name, and last name
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user body dto.CreateUserRequest true "User info"
// @Success      201 {object} dto.SignupResponse "Created user"
// @Failure      400 {object} map[string]string "Bad request / validation error"
// @Failure      409 {object} map[string]string "Email already exists"
// @Failure      429 {object} map[string]string "Too many requests"
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

	user, err := h.service.Register(req.Email, req.Password, req.FirstName, req.LastName)
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
// @Success      200 {object} dto.LoginUserResponse "Login successful, returns tokens"
// @Failure      400 {object} map[string]string "Bad request / validation error"
// @Failure      401 {object} map[string]string "Invalid credentials"
// @Failure      429 {object} map[string]string "Too many requests"
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

	accessToken, refreshToken, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		return err
	}

	claims, _ := auth.ValidateToken(accessToken)

	return c.Status(fiber.StatusOK).JSON(dto.LoginUserResponse{
		UserID:       claims.UserID,
		UserRole:     claims.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// RefreshToken 	godoc
// @Summary      Refresh access token
// @Description  Get a new access token using a refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body dto.RefreshTokenRequest true "Refresh Token"
// @Success      200 {object} dto.RefreshTokenResponse "New tokens"
// @Failure      400 {object} map[string]string "Bad request"
// @Failure      401 {object} map[string]string "Invalid token"
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if ok := utils.ValidateStruct(c, &req); !ok {
		return nil
	}

	accessToken, refreshToken, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Logout 		 godoc
// @Summary      Logout
// @Description  Invalidate access/refresh tokens
// @Tags         Auth
// @Security     Bearer
// @Success      200 {object} map[string]string "Logged out"
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
	}

	var req dto.LogoutRequest
	if err := c.BodyParser(&req); err == nil && req.RefreshToken != "" {
		_ = h.service.InvalidateRefreshToken(req.RefreshToken)
	}

	if err := h.service.Logout("", tokenString); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logged out successfully"})
}
