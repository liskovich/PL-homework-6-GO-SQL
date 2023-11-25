package service

import (
	"example.com/api/db"
	"example.com/api/model"
	"github.com/go-playground/validator/v10"
)

type AuthService interface {
	Login(user model.User)
	Register(user model.UserMutate)
}

type AuthServiceImpl struct {
	UserRepository db.UserRepository
	Validate       *validator.Validate
}

func NewAuthService(userRepository db.UserRepository, validate *validator.Validate) AuthService {
	return &AuthServiceImpl{
		UserRepository: userRepository,
		Validate:       validate,
	}
}

// Login implements AuthService.
func (*AuthServiceImpl) Login(user model.User) {
	panic("unimplemented")
}

// Register implements AuthService.
func (*AuthServiceImpl) Register(user model.UserMutate) {
	panic("unimplemented")
}
