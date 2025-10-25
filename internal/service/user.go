//go:generate mockgen -source=user.go -destination=mock_user.go -package=service -typed
package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token expired")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidTypeToken   = errors.New("invalid token type")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

type UserService interface {
	ValidateCredentials(ctx context.Context, username, password string) (bool, error)
	GenerateToken(ctx context.Context, username string, TokenType string) (string, error)
	RefreshToken(ctx context.Context, token string) (string, error)
}

type User struct {
	Username string
	Password string
}

type UserStore interface {
	Get(ctx context.Context, username string) (*User, error)
}

type userService struct {
	store     UserStore
	jwtSecret string
}

func NewUserService(store UserStore, jwtSecret string) UserService {
	return &userService{
		store:     store,
		jwtSecret: jwtSecret,
	}
}

func (s *userService) ValidateCredentials(ctx context.Context, username, password string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	user, err := s.store.Get(ctx, username)
	if err != nil {
		return false, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, ErrInvalidCredentials
	}

	return true, nil
}

func (s *userService) GenerateToken(ctx context.Context, username string, tokenType string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	var ttl time.Duration
	switch tokenType {
	case TokenTypeAccess:
		ttl = time.Hour
	case TokenTypeRefresh:
		ttl = 30 * 24 * time.Hour
	default:
		return "", ErrInvalidTypeToken
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  username,
		"exp":  time.Now().Add(ttl).Unix(),
		"type": tokenType,
	})

	return token.SignedString([]byte(s.jwtSecret))
}

func (s *userService) RefreshToken(ctx context.Context, tokenString string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return "", ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Проверка на refresh токен
		if tokenType, ok := claims["type"].(string); !ok || tokenType != TokenTypeRefresh {
			return "", ErrInvalidTypeToken
		}

		if username, ok := claims["sub"].(string); ok {
			return s.GenerateToken(ctx, username, TokenTypeAccess)
		}
	}

	return "", ErrInvalidToken
}
