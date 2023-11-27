package service

import (
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
	if err != nil {
		panic(err)
	}
	if user == nil {
		panic("user not found")
	}
	// Compare password with hashed password
	pwdErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if pwdErr != nil {
		panic("invalid password")
	}
	return user
}

func (authSrvc *AuthServiceImpl) Register(user model.UserMutate) {
	existingUser, err := authSrvc.UserRepository.GetUserByEmail(user.Email)
	if err != nil {
		panic(err)
	}
	if existingUser != nil {
		panic("user already registered")
	}
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
}

func (authSrvc *AuthServiceImpl) GetUserByID(userID uint) *model.User {
	user, err := authSrvc.UserRepository.GetUserById(userID)
	if err != nil {
		panic(err)
	}
	return user
}
