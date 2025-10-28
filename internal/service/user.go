//go:generate mockgen -source=user.go -destination=mock_user.go -package=service -typed
package service

import (
	"auth_test/internal/store"
	"auth_test/pkg/metrics"
	"context"
	"errors"
	"log"
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
	CreateUser(ctx context.Context, username, password string) error
}

type User struct {
	ID       int
	Username string
	Password string
}

type UserStore interface {
	Get(ctx context.Context, username string) (*User, error)
	Store(ctx context.Context, user User) error
}

type userService struct {
	store     *store.PostgresStore
	jwtSecret string
}

func NewUserService(store *store.PostgresStore, jwtSecret string) UserService {
	return &userService{
		store:     store,
		jwtSecret: jwtSecret,
	}
}

func (s *userService) ValidateCredentials(ctx context.Context, username, password string) (bool, error) {
	start := time.Now()
	defer func() {
		metrics.LoginDuration.Observe(time.Since(start).Seconds())
	}()

	if err := ctx.Err(); err != nil {
		metrics.LoginAttempts.WithLabelValues("failure").Inc()
		return false, err
	}

	user, err := s.store.GetUser(ctx, username)
	if err != nil {
		metrics.LoginAttempts.WithLabelValues("failure").Inc()
		return false, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		metrics.LoginAttempts.WithLabelValues("failure").Inc()
		return false, ErrInvalidCredentials
	}

	metrics.LoginAttempts.WithLabelValues("success").Inc()
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

	metrics.TokenGenerated.WithLabelValues(tokenType).Inc()

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
		metrics.TokensValidated.WithLabelValues("invalid").Inc()
		return "", ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Проверка на refresh токен
		if tokenType, ok := claims["type"].(string); !ok || tokenType != TokenTypeRefresh {
			metrics.TokensValidated.WithLabelValues("invalid_type").Inc()
			return "", ErrInvalidTypeToken
		}

		if username, ok := claims["sub"].(string); ok {
			metrics.TokensValidated.WithLabelValues("valid").Inc()
			return s.GenerateToken(ctx, username, TokenTypeAccess)
		}
	}

	return "", ErrInvalidToken
}

func (s *userService) CreateUser(ctx context.Context, username, password string) error {
	start := time.Now()

	if err := ctx.Err(); err != nil {
		metrics.UserCreated.WithLabelValues("failure").Inc()
		return err
	}

	log.Printf("Creating user: %s", username)

	exists, err := s.store.UserExists(ctx, username)
	if err != nil {
		return err
	}
	if exists {
		return ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to create user %s: %v", username, err)
		metrics.UserCreated.WithLabelValues("failure").Inc()
		return err
	}

	user := store.CreateUserParams{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	_, err = s.store.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			metrics.UserCreated.WithLabelValues("conflict").Inc()
		} else {
			metrics.UserCreated.WithLabelValues("failure").Inc()
		}
		log.Printf("Failed to create user %s: %v", username, err)
		return err
	}

	metrics.UserCreated.WithLabelValues("success").Inc()
	metrics.UserCreationDuration.Observe(time.Since(start).Seconds())
	log.Printf("User created successfully: %s", username)
	return nil
}
