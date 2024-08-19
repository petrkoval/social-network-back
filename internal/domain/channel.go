package domain

import "time"

type Channel struct {
	ID          string    `json:"id" db:"channel_id"`
	UserID      string    `json:"user_id" db:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
}

type CreateChannelDTO struct {
	UserID      string `json:"user_id" db:"user_id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
}

type UpdateChannelDTO struct {
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
}
