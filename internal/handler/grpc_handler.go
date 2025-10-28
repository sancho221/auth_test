package handler

import (
	"auth_test/internal/service"
	"auth_test/pkg/pb"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCHandler struct {
	pb.UnimplementedAuthServiceServer
	userService service.UserService
}

func NewGRPCHandler(userService service.UserService) *GRPCHandler {
	return &GRPCHandler{
		userService: userService,
	}
}

func (h *GRPCHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	valid, err := h.userService.ValidateCredentials(ctx, req.Username, req.Password)
	if err != nil || !valid {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	accessToken, err := h.userService.GenerateToken(ctx, req.Username, service.TokenTypeAccess)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate access token")
	}
	refreshToken, err := h.userService.GenerateToken(ctx, req.Username, service.TokenTypeRefresh)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate refresh token")
	}

	return &pb.LoginResponse{
		Message:      "Login successful",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}, nil
}

func (h *GRPCHandler) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if req.Token == "" {
		return &pb.VerifyTokenResponse{
			Message: "Token is invalid",
			Valid:   false,
		}, nil
	}

	newAccessToken, err := h.userService.RefreshToken(ctx, req.Token)
	if err != nil {
		return &pb.VerifyTokenResponse{
			Message: "Token is invalid",
			Valid:   false,
		}, nil
	}

	if newAccessToken != req.Token {
		return &pb.VerifyTokenResponse{
			Message:     "Token refreshed",
			AccessToken: newAccessToken,
			TokenType:   "Bearer",
			Valid:       true,
		}, nil
	}

	return &pb.VerifyTokenResponse{
		Message: "Token is valid",
		Valid:   true,
	}, nil
}
