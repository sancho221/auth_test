package handler

import (
	"auth_test/internal/service"
	"auth_test/pkg/pb"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthService_Login(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		password      string
		mockValid     bool
		mockErr       error
		expectedCode  codes.Code
		expectedToken string
		expectedError string
	}{
		{
			name:          "successful login",
			username:      "admin",
			password:      "admin123",
			mockValid:     true,
			mockErr:       nil,
			expectedCode:  codes.OK,
			expectedToken: "access-token-123",
		},
		{
			name:          "invalid credentials - wrong password",
			username:      "admin",
			password:      "wrongpassword",
			mockValid:     false,
			mockErr:       service.ErrInvalidCredentials,
			expectedCode:  codes.Unauthenticated,
			expectedError: "invalid credentials",
		},
		{
			name:          "invalid credentials - user not found",
			username:      "nonexists",
			password:      "anypassword",
			mockValid:     false,
			mockErr:       service.ErrUserNotFound,
			expectedCode:  codes.Unauthenticated,
			expectedError: "invalid credentials",
		},
		{
			name:          "server error",
			username:      "admin",
			password:      "admin123",
			mockValid:     false,
			mockErr:       errors.New("database connection failed"),
			expectedCode:  codes.Unauthenticated,
			expectedError: "invalid credentials",
		},
		{
			name:          "empty credentials",
			username:      "",
			password:      "",
			mockValid:     false,
			mockErr:       nil,
			expectedCode:  codes.Unauthenticated,
			expectedError: "invalid credentials",
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
			} else {
				mockService.EXPECT().ValidateCredentials(gomock.Any(), "", "").Return(false, nil)
			}

			handler := NewGRPCHandler(mockService)

			req := &pb.LoginRequest{
				Username: tt.username,
				Password: tt.password,
			}

			resp, err := handler.Login(context.Background(), req)

			if tt.expectedCode == codes.OK {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, tt.expectedToken, resp.AccessToken)
				assert.Equal(t, "refresh-token-456", resp.RefreshToken)
				assert.Equal(t, "Bearer", resp.TokenType)
			} else {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tt.expectedCode, st.Code())
				assert.Contains(t, st.Message(), tt.expectedError)
			}
		})
	}
}
