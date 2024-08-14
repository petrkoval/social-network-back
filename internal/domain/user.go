package domain

type User struct {
	ID       string `json:"id"       db:"user_id"`
	Username string `json:"username" db:"username"`
}
