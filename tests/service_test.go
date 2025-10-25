package tests

import (
	"auth_test/internal/service"
	"auth_test/internal/store"
	"testing"
)

// UserService method ValidateCredentials
func TestUserService_ValidateCredentials(t *testing.T) {
	userStore := store.NewInMemoryStore()
	userService := service.NewUserService(userStore, "test-secret")

	valid, err := userService.ValidateCredentials("admin", "admin123")
	if err != nil || !valid {
		t.Errorf("Excepted valid credentials, got error: %v", err)
	}

	valid, err = userService.ValidateCredentials("admin", "wrongpassword")
	if err == nil || valid {
		t.Error("Excepted invalid credentials error")
	}
}

// UserService method GenerateToken
func TestUserService_GenerateToken(t *testing.T) {
	userStore := store.NewInMemoryStore()
	userService := service.NewUserService(userStore, "test-secret")

	token, err := userService.GenerateToken("admin")
	if err != nil {
		t.Errorf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Error("Token should not be empty")
	}
}

// UserService method RefreshToken
func TestUserService_RefreshToken(t *testing.T) {
	userStore := store.NewInMemoryStore()
	userService := service.NewUserService(userStore, "test-secret")

	token, err := userService.GenerateToken("admin")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	newToken, err := userService.RefreshToken(token)
	if err != nil {
		t.Errorf("Refresh failed: %v", err)
	}
	if newToken == "" {
		t.Error("New token should not be empty")
	}

	_, err = userService.RefreshToken("invalid-token")
	if err == nil {
		t.Error("Excepted error for invalid token")
	}
}

// UserService проверка ValidateCredentials на неизвестного пользователяя
func TestUserService_ValidateCredentials_UserNotFound(t *testing.T) {
	userStore := store.NewInMemoryStore()
	userService := service.NewUserService(userStore, "test-secret")

	valid, err := userService.ValidateCredentials("empty", "password")
	if err == nil || valid {
		t.Error("Expected error for non-existent user")
	}
}

// UserService проверка RefreshToken на просроченный токен
func TestUserService_RefreshToken_ExpiredToken(t *testing.T) {
	userStore := store.NewInMemoryStore()
	userService := service.NewUserService(userStore, "test-secret")

	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhZG1pbiIsImV4cCI6MTUxNjIzOTAyMn0.invalid-signature"

	_, err := userService.RefreshToken(expiredToken)
	if err == nil {
		t.Error("Excepted error for expired/invalid token")
	}
}
