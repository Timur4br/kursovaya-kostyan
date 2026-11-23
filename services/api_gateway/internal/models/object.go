package models

import (
	"time"
)

type Object struct {
	FirstSeen    *time.Time `json:"first_seen"`
	ID           int64      `json:"id"`
	ImagePath    *string    `json:"image_path"`
	LastSeen     *time.Time `json:"last_seen"`
	Name         string     `json:"name"`
	ObjectTypeID int64      `json:"object_type_id"`
	Status       string     `json:"status"`
}
