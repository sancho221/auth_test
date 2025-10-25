package store

import (
	"auth_test/internal/service"
	"context"

	"golang.org/x/crypto/bcrypt"
)

type InMemoryStore struct {
	users map[string]service.User
}

func NewInMemoryStore() *InMemoryStore {
	adminHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	userHash, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)

	return &InMemoryStore{
		users: map[string]service.User{
			"admin": {Username: "admin", Password: string(adminHash)},
			"user":  {Username: "user", Password: string(userHash)},
		},
	}
}

func (s *InMemoryStore) Store(ctx context.Context, username, password string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, exists := s.users[username]; exists {
		return service.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	s.users[username] = service.User{
		Username: username,
		Password: string(hashedPassword),
	}
	return nil
}

func (s *InMemoryStore) Get(ctx context.Context, username string) (*service.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	user, exists := s.users[username]
	if !exists {
		return nil, service.ErrUserNotFound
	}

	return &user, nil
}
