package tests

type MockUserService struct {
	ValidateCredentialsFunc func(username, password string) (bool, error)
	GenerateTokenFunc       func(username string) (string, error)
	RefreshTokenFunc        func(token string) (string, error)
}

func (m *MockUserService) ValidateCredentials(username, password string) (bool, error) {
	return m.ValidateCredentialsFunc(username, password)
}

func (m *MockUserService) GenerateToken(username string) (string, error) {
	return m.GenerateTokenFunc(username)
}

func (m *MockUserService) RefreshToken(token string) (string, error) {
	return m.RefreshTokenFunc(token)
}
