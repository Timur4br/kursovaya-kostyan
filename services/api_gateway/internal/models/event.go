package models

import (
	"time"
)

type Event struct {
	CameraID    int64     `json:"camera_id"`
	Confidence  *float64  `json:"confidence"`
	Description string    `json:"description"`
	EventTypeID int64     `json:"event_type_id"`
	ID          int64     `json:"id"`
	LocationID  int64     `json:"location_id"`
	ObjectID    *int64    `json:"object_id"`
	OccurredAt  time.Time `json:"occurred_at"`
	PersonID    *int64    `json:"person_id"`
}
