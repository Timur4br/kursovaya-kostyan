package entities

import "time"

type Object struct {
	ID        int       `json:"id"`
	Type      string    `json:"type" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
