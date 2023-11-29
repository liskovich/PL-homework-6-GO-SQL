package service

import (
	"database/sql"

	"example.com/api/db"
	"example.com/api/model"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(email string, password string) *model.User
	Register(user model.UserMutate)
	GetUserByID(userID uint) *model.User
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

func (authSrvc *AuthServiceImpl) Login(email string, password string) *model.User {
	user, err := authSrvc.UserRepository.GetUserByEmail(email)
	switch {
	case err == sql.ErrNoRows:
		panic("invalid email or password")
	case err != nil:
		panic(err)
	default:
		// Compare password with hashed password
		pwdErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if pwdErr != nil {
			panic("invalid password")
		}
		return user
	}
}

func (authSrvc *AuthServiceImpl) Register(user model.UserMutate) {
	_, err := authSrvc.UserRepository.GetUserByEmail(user.Email)
	switch {
	case err == sql.ErrNoRows:
		// Hash password
		hashedPassword, pwdErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if pwdErr != nil {
			panic(pwdErr)
		}
		userToCreate := model.UserMutate{
			Name:     user.Name,
			Email:    user.Email,
			Password: string(hashedPassword),
		}
		valErr := authSrvc.Validate.Struct(userToCreate)
		if valErr != nil {
			panic(valErr)
		}
		saveErr := authSrvc.UserRepository.CreateUser(userToCreate)
		if saveErr != nil {
			panic(saveErr)
		}
	case err != nil:
		panic(err)
	default:
		panic("user already registered")
	}
}

func (authSrvc *AuthServiceImpl) GetUserByID(userID uint) *model.User {
	user, err := authSrvc.UserRepository.GetUserById(userID)
	switch {
	case err == sql.ErrNoRows:
		panic("user not found")
	case err != nil:
		panic(err)
	default:
		return user
	}
}
