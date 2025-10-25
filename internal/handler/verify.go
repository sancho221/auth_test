package handler

import (
	"auth_test/internal/service"
	"net/http"
	"strings"
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
	ctx := r.Context()

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		JSONError(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		JSONError(w, "Invalid Authorization format: expected 'Bearer <token>'", http.StatusUnauthorized)
		return
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		JSONError(w, "Invalid Authorization scheme: expected Bearer", http.StatusUnauthorized)
		return
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		JSONError(w, "Authorization token is empty", http.StatusUnauthorized)
		return
	}

	newAccessToken, err := h.userService.RefreshToken(ctx, token)
	if err != nil {
		JSONError(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	response := VerifyResponse{
		Message:     "Token refreshed",
		AccessToken: newAccessToken,
		TokenType:   "Bearer",
	}

	JSONSuccess(w, response, http.StatusOK)

}
