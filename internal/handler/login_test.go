package handler

import (
	"auth_test/internal/service"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLoginHandler_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		password       string
		mockValid      bool
		mockErr        error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful login",
			username:       "admin",
			password:       "admin123",
			mockValid:      true,
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "access-token-123",
		},
		{
			name:           "invalid credentials - wrong password",
			username:       "admin",
			password:       "wrongpassword",
			mockValid:      false,
			mockErr:        service.ErrInvalidCredentials,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid credentials",
		},
		{
			name:           "invalid credentials - user not found",
			username:       "nonexists",
			password:       "anypassword",
			mockValid:      false,
			mockErr:        service.ErrUserNotFound,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid credentials",
		},
		{
			name:           "server error",
			username:       "admin",
			password:       "admin123",
			mockValid:      false,
			mockErr:        errors.New("database connection failed"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Internal server error",
		},
		{
			name:           "missing basic auth",
			username:       "",
			password:       "",
			mockValid:      false,
			mockErr:        nil,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Missing or invalid Authorization header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := service.NewMockUserService(ctrl)

			if tt.username != "" && tt.password != "" {
				mockService.EXPECT().ValidateCredentials(gomock.Any(), tt.username, tt.password).Return(tt.mockValid, tt.mockErr)

				if tt.mockValid && tt.mockErr == nil {
					mockService.EXPECT().GenerateToken(gomock.Any(), tt.username, service.TokenTypeAccess).Return("access-token-123", nil)
					mockService.EXPECT().GenerateToken(gomock.Any(), tt.username, service.TokenTypeRefresh).Return("refresh-token-456", nil)
				}
			}

			handler := NewLoginHandler(mockService)
			req := httptest.NewRequest("POST", "/login", nil)

			if tt.username != "" && tt.password != "" {
				req.SetBasicAuth(tt.username, tt.password)
			}

			rr := httptest.NewRecorder()
			handler.Handle(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code, "Status code mismatch")
			assert.Contains(t, rr.Body.String(), tt.expectedBody, "Responsy body mismatch")
		})
	}
}
