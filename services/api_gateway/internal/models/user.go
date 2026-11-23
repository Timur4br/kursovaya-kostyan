package models

import (
	"time"
)

type User struct {
	CreatedAt time.Time `json:"created_at"`
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
}
