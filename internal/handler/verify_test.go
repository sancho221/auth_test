package handler

import (
	"auth_test/internal/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestVerifyHandler_TableDriven(t *testing.T) {
	tests := []struct {
		name            string
		authHeader      string
		mockReturnToken string
		mockErr         error
		expectedStatus  int
		expectedBody    string
	}{
		{
			name:            "successful token refresh",
			authHeader:      "Bearer valid-token",
			mockReturnToken: "new-access-token",
			mockErr:         nil,
			expectedStatus:  http.StatusOK,
			expectedBody:    "new-access-token",
		},
		{
			name:            "expired token",
			authHeader:      "Bearer expired-token",
			mockReturnToken: "",
			mockErr:         service.ErrExpiredToken,
			expectedStatus:  http.StatusUnauthorized,
			expectedBody:    "Invalid or expired token",
		},
		{
			name:            "invalid token",
			authHeader:      "Bearer invalid-token",
			mockReturnToken: "",
			mockErr:         service.ErrInvalidToken,
			expectedStatus:  http.StatusUnauthorized,
			expectedBody:    "Invalid or expired token",
		},
		{
			name:            "missing authorization header",
			authHeader:      "",
			mockReturnToken: "",
			mockErr:         nil,
			expectedStatus:  http.StatusUnauthorized,
			expectedBody:    "Missing Authorization header",
		},
		{
			name:            "invalid authorization format",
			authHeader:      "Bearer token extra",
			mockReturnToken: "",
			mockErr:         nil,
			expectedStatus:  http.StatusUnauthorized,
			expectedBody:    "Invalid Authorization format",
		},
		{
			name:            "wrong authorization scheme",
			authHeader:      "Basic token123",
			mockReturnToken: "",
			mockErr:         nil,
			expectedStatus:  http.StatusUnauthorized,
			expectedBody:    "Invalid Authorization scheme",
		},
		{
			name:            "empty token",
			authHeader:      "Bearer ",
			mockReturnToken: "",
			mockErr:         nil,
			expectedStatus:  http.StatusUnauthorized,
			expectedBody:    "Authorization token is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := service.NewMockUserService(ctrl)

			if tt.authHeader != "" {
				parts := strings.Split(tt.authHeader, " ")
				if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") && strings.TrimSpace(parts[1]) != "" {
					token := strings.TrimSpace(parts[1])
					mockService.EXPECT().RefreshToken(gomock.Any(), token).Return(tt.mockReturnToken, tt.mockErr)
				}
			}

			handler := NewVerifyHandler(mockService)
			req := httptest.NewRequest("POST", "/verify", nil)

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()
			handler.Handle(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}
