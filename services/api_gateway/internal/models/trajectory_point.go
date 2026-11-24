package models

import (
	"time"
)

type TrajectoryPoint struct {
	ID             int64     `json:"id"`
	ObservedAt     time.Time `json:"observed_at"`
	TrajectoryID   int64     `json:"trajectory_id"`
	VideoSegmentID *int64    `json:"video_segment_id"`
	X              float64   `json:"x"`
	Y              float64   `json:"y"`
	Z              *float64  `json:"z"`
}
