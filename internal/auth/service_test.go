package auth_test

import (
	"adv-go/api/internal/auth"
	"adv-go/api/internal/user"
	"testing"
)

const (
	emailWant = "a@a1.ru"
)

type mockUserRepository struct {}

func (mock *mockUserRepository) Create(user *user.User) (*user.User, error) {
	return user, nil
}

func (mock *mockUserRepository) FindByEmail(email string) (*user.User, error) {
	return nil, nil
}

func TestRegisterSuccess(t *testing.T) {
	authService := auth.NewAuthService(&mockUserRepository{})
	emailGot, err := authService.Register(emailWant, "1234", "test user")
	if err != nil {
		t.Fatal(err)
	}
	if emailGot != emailWant {
		t.Fatalf("error during registration. email want: %s, got: %s", emailGot, emailWant)
	}
}