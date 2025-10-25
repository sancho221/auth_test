package handler

import (
	"auth_test/internal/service"
	"encoding/json"
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
	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}

	valid, err := h.userService.ValidateCredentials(username, password)
	if err != nil || !valid {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := h.userService.GenerateToken(username)
	if err != nil {
		println("Token generation error:", err.Error())
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Login successful",
		"token":   token,
	}
	json.NewEncoder(w).Encode(response)
}
