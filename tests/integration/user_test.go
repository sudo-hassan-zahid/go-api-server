package integration

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/sudo-hassan-zahid/go-api-server/internal/dto"
	"github.com/sudo-hassan-zahid/go-api-server/internal/models"
	"github.com/sudo-hassan-zahid/go-api-server/tests"
)

func TestUserFlow(t *testing.T) {
	cleanup, err := tests.SetupTestContainer()
	if err != nil {
		t.Fatalf("setup failed: %s", err)
	}
	defer cleanup()

	if err := tests.TestDB.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("migration failed: %s", err)
	}

	app := tests.SetupApp()

	adminEmail := "admin@example.com"
	adminPass := "password123"

	signupAdminBody, _ := json.Marshal(dto.CreateUserRequest{
		Email:     adminEmail,
		Password:  adminPass,
		FirstName: "Admin",
		LastName:  "User",
	})
	req := httptest.NewRequest("POST", "/api/auth/signup", bytes.NewReader(signupAdminBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	var adminSignupResp dto.SignupResponse
	json.NewDecoder(resp.Body).Decode(&adminSignupResp)

	tests.TestDB.Model(&models.User{}).Where("email = ?", adminEmail).Update("role", "admin")

	var adminLoginResp dto.LoginUserResponse
	loginAdminBody, _ := json.Marshal(dto.LoginUserRequest{
		Email:    adminEmail,
		Password: adminPass,
	})
	req = httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(loginAdminBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	json.NewDecoder(resp.Body).Decode(&adminLoginResp)
	adminToken := adminLoginResp.AccessToken

	userEmail := "user@example.com"
	userPass := "password123"
	signupUserBody, _ := json.Marshal(dto.CreateUserRequest{
		Email:     userEmail,
		Password:  userPass,
		FirstName: "Regular",
		LastName:  "User",
	})
	req = httptest.NewRequest("POST", "/api/auth/signup", bytes.NewReader(signupUserBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	var userSignupResp dto.SignupResponse
	json.NewDecoder(resp.Body).Decode(&userSignupResp)
	userToken := userSignupResp.AccessToken
	userID := userSignupResp.UserID

	req = httptest.NewRequest("GET", "/api/users/"+userID, nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	updateBody, _ := json.Marshal(dto.UpdateUserRequest{
		FirstName: "Updated",
	})
	req = httptest.NewRequest("PATCH", "/api/users/"+userID, bytes.NewReader(updateBody))
	req.Header.Set("Authorization", "Bearer "+userToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var updatedUser models.User
	tests.TestDB.First(&updatedUser, "id = ?", userID)
	assert.Equal(t, "Updated", updatedUser.FirstName)

	req = httptest.NewRequest("DELETE", "/api/users/"+userID, nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)

	req = httptest.NewRequest("DELETE", "/api/users/"+userID, nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	req = httptest.NewRequest("GET", "/api/users/"+userID, nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}
