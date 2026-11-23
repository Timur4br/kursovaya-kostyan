package models

import (
	"time"
)

type VideoSegment struct {
	CameraID         int64     `json:"camera_id"`
	DurationSeconds  float64   `json:"duration_seconds"`
	EndTime          time.Time `json:"end_time"`
	EventID          *int64    `json:"event_id"`
	ID               int64     `json:"id"`
	Preview          *string   `json:"preview"`
	StartTime        time.Time `json:"start_time"`
	StorageReference string    `json:"storage_reference"`
}
