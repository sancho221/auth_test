package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token expired")
	ErrUserNotFound       = errors.New("user not found")
)

type UserService interface {
	ValidateCredentials(username, password string) (bool, error)
	GenerateToken(username string) (string, error)
	RefreshToken(token string) (string, error)
}

type User struct {
	Username string
	Password string
}

type UserStore interface {
	Get(username string) (*User, error)
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

func (s *userService) ValidateCredentials(username, password string) (bool, error) {
	user, err := s.store.Get(username)
	if err != nil {
		return false, ErrInvalidCredentials
	}

	if user.Password != password {
		return false, ErrInvalidCredentials
	}

	return true, nil
}

func (s *userService) GenerateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}

func (s *userService) RefreshToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return "", ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if username, ok := claims["sub"].(string); ok {
			return s.GenerateToken(username)
		}
	}

	return "", ErrInvalidToken
}
