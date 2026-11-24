package models

import (
	"time"
)

type QueryHistory struct {
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	ID          int64      `json:"id"`
	Response    *string    `json:"response"`
	Status      string     `json:"status"`
	Text        string     `json:"text"`
	UserID      int64      `json:"user_id"`
}
