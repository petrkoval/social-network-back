package domain

import "time"

type User struct {
	ID                 string    `json:"id"       db:"user_id"`
	Username           string    `json:"username" db:"username"`
	Password           string    `json:"-" db:"password"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	AccountDescription string    `json:"account_description" db:"account_description"`
}

type CreateUserDTO struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type AuthUser struct {
	ID       string `json:"id"       db:"user_id"`
	Username string `json:"username" db:"username"`
}
