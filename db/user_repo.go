package db

import (
	"database/sql"

	"example.com/api/model"
)

type UserRepository interface {
	CreateUser(user model.UserMutate) error
	UpdateUser(userID uint, user model.User) (*model.User, error)
	GetUserById(userID uint) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (usrRepo *userRepo) CreateUser(user model.UserMutate) error {
	_, err := usrRepo.db.Exec(InsertUserQuery, user.Name, user.Email, user.Password)
	return err
}

func (usrRepo *userRepo) GetUserById(userID uint) (*model.User, error) {
	row := usrRepo.db.QueryRow(SelectUserByIdQuery, userID)
	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (usrRepo *userRepo) UpdateUser(userID uint, user model.User) (*model.User, error) {
	_, err := usrRepo.db.Exec(UpdateUserQuery, user.Name, user.Email, user.Password, userID)
	if err != nil {
		return nil, err
	}

	updatedUser, err := usrRepo.GetUserById(userID)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (usrRepo *userRepo) GetUserByEmail(email string) (*model.User, error) {
	row := usrRepo.db.QueryRow(SelectUserByEmailQuery, email)
	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
