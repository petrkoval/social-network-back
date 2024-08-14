package domain

type User struct {
	ID       string `json:"id"       db:"user_id"`
	Username string `json:"username" db:"username"`
	Password string `json:"-" db:"password"`
}

type CreateUserDTO struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}
