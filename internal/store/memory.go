package store

import (
	"auth_test/internal/service"
)

var ()

type InMemoryStore struct {
	users map[string]service.User
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		users: map[string]service.User{
			"admin": {Username: "admin", Password: "admin123"},
			"user":  {Username: "user", Password: "user123"},
		},
	}
}

func (s *InMemoryStore) Store(username, password string) {
	// здесь по хорошему нужно сделать проверку (на существование логина) + тесты
	// ...

	s.users[username] = service.User{
		Username: username,
		Password: password,
	}
}

func (s *InMemoryStore) Get(username string) (*service.User, error) {
	user, exists := s.users[username]
	if !exists {
		return nil, service.ErrUserNotFound
	}

	return &user, nil
}
