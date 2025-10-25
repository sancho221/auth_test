package tests

import (
	"auth_test/internal/handler"
	"auth_test/internal/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Успешный логин
func TestLoginHandler_Success(t *testing.T) {
	mockService := &MockUserService{
		ValidateCredentialsFunc: func(username, password string) (bool, error) {
			return true, nil
		},
		GenerateTokenFunc: func(username string) (string, error) {
			return "test-token", nil
		},
	}

	handler := handler.NewLoginHandler(mockService)
	req := httptest.NewRequest("POST", "/login", nil)
	req.SetBasicAuth("admin", "admin123")

	rr := httptest.NewRecorder()
	handler.Handle(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Excepted status 200, got %d", status)
	}

	authHeader := rr.Header().Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		t.Error("Excepted Authorization header with Bearer token")
	}
}

// Неверный пароль
func TestLoginHandler_InvalidCredentials(t *testing.T) {
	mockService := &MockUserService{
		ValidateCredentialsFunc: func(username, password string) (bool, error) {
			return false, service.ErrInvalidCredentials
		},
	}

	handler := handler.NewLoginHandler(mockService)
	req := httptest.NewRequest("POST", "/login", nil)
	req.SetBasicAuth("admin", "wrongpassword")

	rr := httptest.NewRecorder()
	handler.Handle(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Excepted status 401, got %d", status)
	}
}

// Отсутствие токена
func TestLoginHandler_MissingAuthHeader(t *testing.T) {
	mockService := &MockUserService{}

	handler := handler.NewLoginHandler(mockService)
	req := httptest.NewRequest("POST", "/login", nil)

	rr := httptest.NewRecorder()
	handler.Handle(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Excepted status 401, got %d", status)
	}
}
