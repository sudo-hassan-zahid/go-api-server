package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/sudo-hassan-zahid/go-api-server/internal/dto"
	"github.com/sudo-hassan-zahid/go-api-server/internal/models"
	"github.com/sudo-hassan-zahid/go-api-server/tests"
)

func TestAuthFlow(t *testing.T) {
	cleanup, err := tests.SetupTestContainer()
	if err != nil {
		t.Fatalf("setup failed: %s", err)
	}
	defer cleanup()

	if err := tests.TestDB.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("migration failed: %s", err)
	}

	app := tests.SetupApp()

	t.Run("Full Auth Lifecycle", func(t *testing.T) {
		email := "e2e@example.com"
		password := "password123"

		signupBody, _ := json.Marshal(dto.CreateUserRequest{
			Email:     email,
			Password:  password,
			FirstName: "E2E",
			LastName:  "Test",
		})
		req := httptest.NewRequest("POST", "/api/auth/signup", bytes.NewReader(signupBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		loginBody, _ := json.Marshal(dto.LoginUserRequest{
			Email:    email,
			Password: password,
		})
		req = httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(loginBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err = app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var loginResp dto.LoginUserResponse
		json.NewDecoder(resp.Body).Decode(&loginResp)
		assert.NotEmpty(t, loginResp.AccessToken)
		assert.NotEmpty(t, loginResp.RefreshToken)

		refreshBody, _ := json.Marshal(dto.RefreshTokenRequest{
			RefreshToken: loginResp.RefreshToken,
		})
		req = httptest.NewRequest("POST", "/api/auth/refresh", bytes.NewReader(refreshBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err = app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var refreshResp dto.RefreshTokenResponse
		json.NewDecoder(resp.Body).Decode(&refreshResp)
		assert.NotEmpty(t, refreshResp.AccessToken)
		assert.NotEmpty(t, refreshResp.RefreshToken)
		assert.NotEqual(t, loginResp.AccessToken, refreshResp.AccessToken)
		assert.NotEqual(t, loginResp.RefreshToken, refreshResp.RefreshToken)

		req = httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(loginBody))
		req.Header.Set("Content-Type", "application/json")
		for i := 0; i < 6; i++ {
			resp, _ = app.Test(req)
		}

		key := "login_attempts:"
		keys, _ := tests.TestRdb.Keys(context.Background(), key+"*").Result()
		if len(keys) > 0 {
			count, _ := tests.TestRdb.Get(context.Background(), keys[0]).Int()
			assert.GreaterOrEqual(t, count, 6)
		}

		logoutBody, _ := json.Marshal(dto.LogoutRequest{
			RefreshToken: refreshResp.RefreshToken,
		})
		req = httptest.NewRequest("POST", "/api/auth/logout", bytes.NewReader(logoutBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+refreshResp.AccessToken)
		resp, err = app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		req = httptest.NewRequest("POST", "/api/auth/refresh", bytes.NewReader(refreshBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err = app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})
}
