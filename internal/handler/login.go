package handler

import (
	"auth_test/internal/service"
	"errors"
	"net/http"
)

type LoginHandler struct {
	userService service.UserService
}

func NewLoginHandler(userService service.UserService) *LoginHandler {
	return &LoginHandler{
		userService: userService,
	}
}

func (h *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username, password, ok := r.BasicAuth()
	if !ok {
		JSONError(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}

	valid, err := h.userService.ValidateCredentials(ctx, username, password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) || errors.Is(err, service.ErrUserNotFound) {
			JSONError(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			JSONError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	if !valid {
		JSONError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, err := h.userService.GenerateToken(ctx, username, service.TokenTypeAccess)
	if err != nil {
		JSONError(w, "Failed to generate assecc token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := h.userService.GenerateToken(ctx, username, service.TokenTypeRefresh)
	if err != nil {
		JSONError(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		Message:      "Login successful",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}

	JSONSuccess(w, response, http.StatusOK)
}
