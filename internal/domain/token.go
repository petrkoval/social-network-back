package domain

type Token struct {
	UserID       string `json:"user_id"       db:"user_id"`
	RefreshToken string `json:"refresh_token" db:"refresh_token"`
}
