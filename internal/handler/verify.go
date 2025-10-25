package handler

import (
	"auth_test/internal/service"
	"encoding/json"
	"net/http"
)

type VerifyHandler struct {
	userService service.UserService
}

func NewVerifyHandler(userService service.UserService) *VerifyHandler {
	return &VerifyHandler{
		userService: userService,
	}
}

func (h *VerifyHandler) Handle(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		http.Error(w, "Invalid Authorization format", http.StatusUnauthorized)
		return
	}

	token := authHeader[7:]
	newToken, err := h.userService.RefreshToken(token)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Token refreshed",
		"token":   newToken,
	}
	json.NewEncoder(w).Encode(response)

}
