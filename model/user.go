package model

type User struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password_hash"`
}

type UserMutate struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password_hash"`
}
