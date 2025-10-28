package handler

import (
	"auth_test/internal/service"
	"auth_test/pkg/pb"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthService_VerifyToken(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		mockReturnToken string
		mockErr         error
		expectedValid   bool
		expectedToken   string
		expectedMessage string
	}{
		{
			name:            "successful token refresh",
			token:           "valid-token",
			mockReturnToken: "new-access-token",
			mockErr:         nil,
			expectedValid:   true,
			expectedToken:   "new-access-token",
			expectedMessage: "Token refreshed",
		},
		{
			name:            "token is already valid",
			token:           "valid-token",
			mockReturnToken: "valid-token",
			mockErr:         nil,
			expectedValid:   true,
			expectedToken:   "",
			expectedMessage: "Token is valid",
		},
		{
			name:            "expired token",
			token:           "expired-token",
			mockReturnToken: "",
			mockErr:         service.ErrExpiredToken,
			expectedValid:   false,
			expectedToken:   "",
			expectedMessage: "Token is invalid",
		},
		{
			name:            "invalid token",
			token:           "invalid-token",
			mockReturnToken: "",
			mockErr:         service.ErrInvalidToken,
			expectedValid:   false,
			expectedToken:   "",
			expectedMessage: "Token is invalid",
		},
		{
			name:            "empty token",
			token:           "",
			mockReturnToken: "",
			mockErr:         nil,
			expectedValid:   false,
			expectedToken:   "",
			expectedMessage: "Token is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := service.NewMockUserService(ctrl)

			if tt.token != "" {
				mockService.EXPECT().RefreshToken(gomock.Any(), tt.token).Return(tt.mockReturnToken, tt.mockErr)
			}

			handler := NewGRPCHandler(mockService)

			req := &pb.VerifyTokenRequest{
				Token: tt.token,
			}

			resp, err := handler.VerifyToken(context.Background(), req)

			require.NoError(t, err)
			require.NotNil(t, resp)

			assert.Equal(t, tt.expectedValid, resp.Valid)
			assert.Equal(t, tt.expectedMessage, resp.Message)

			if tt.expectedToken != "" {
				assert.Equal(t, tt.expectedToken, resp.AccessToken)
				assert.Equal(t, "Bearer", resp.TokenType)
			} else {
				assert.Empty(t, resp.AccessToken)
				assert.Empty(t, resp.TokenType)
			}
		})
	}
}
