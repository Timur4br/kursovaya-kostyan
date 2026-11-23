package models

import (
	"time"
)

type Trajectory struct {
	Direction       *string    `json:"direction"`
	EndLocationID   *int64     `json:"end_location_id"`
	EndTime         *time.Time `json:"end_time"`
	ID              int64      `json:"id"`
	PersonID        int64      `json:"person_id"`
	StartLocationID *int64     `json:"start_location_id"`
	StartTime       time.Time  `json:"start_time"`
	Status          string     `json:"status"`
}
