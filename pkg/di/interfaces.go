package di

import (
	"adv-go/api/internal/user"
)

type IStatRepository interface {
	AddClick(linkId uint)
}

type IUserRepository interface {
	Create(user *user.User) (*user.User, error)
	FindByEmail(email string) (*user.User, error)
}