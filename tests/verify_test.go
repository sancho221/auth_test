package tests

import (
	"auth_test/internal/handler"
	"auth_test/internal/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Успешный токен
func TestVerifyHandler_ValidToken(t *testing.T) {
	mockService := &MockUserService{
		RefreshTokenFunc: func(token string) (string, error) {
			return "new-token", nil
		},
	}

	handler := handler.NewVerifyHandler(mockService)
	req := httptest.NewRequest("POST", "/verify", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	rr := httptest.NewRecorder()
	handler.Handle(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Excepted status 200, got %d", status)
	}

	authHeader := rr.Header().Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		t.Error("Excepted Authorization header with new Bearer token")
	}
}

// Просроченный токен
func TestVerifyHandler_ExpiredToken(t *testing.T) {
	mockService := &MockUserService{
		RefreshTokenFunc: func(token string) (string, error) {
			return "", service.ErrExpiredToken
		},
	}

	handler := handler.NewVerifyHandler(mockService)
	req := httptest.NewRequest("POST", "/verify", nil)
	req.Header.Set("Authorization", "Bearer expired-token")

	rr := httptest.NewRecorder()
	handler.Handle(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Excepted status 401, got %d", status)
	}
}

// Неверный токен
func TestVerifyHandler_InvalidToken(t *testing.T) {
	mockService := &MockUserService{
		RefreshTokenFunc: func(token string) (string, error) {
			return "", service.ErrInvalidToken
		},
	}

	handler := handler.NewVerifyHandler(mockService)
	req := httptest.NewRequest("POST", "/verify", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	rr := httptest.NewRecorder()
	handler.Handle(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Excepted status 401, got %d", status)
	}
}

// Пустой токен
func TestVerifyHandler_MissingAuthHeader(t *testing.T) {
	mockService := &MockUserService{}

	handler := handler.NewVerifyHandler(mockService)
	req := httptest.NewRequest("POST", "/verify", nil)

	rr := httptest.NewRecorder()
	handler.Handle(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Excepted status 401, got %d", status)
	}
}

// Другой (неверный) формат токена
func TestVerifyHandler_InvalidAuthFormat(t *testing.T) {
	mockService := &MockUserService{}

	handler := handler.NewVerifyHandler(mockService)
	req := httptest.NewRequest("POST", "/verify", nil)
	req.Header.Set("Authorization", "Basic token")

	rr := httptest.NewRecorder()
	handler.Handle(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Excepted status 401, got %d", status)
	}
}
