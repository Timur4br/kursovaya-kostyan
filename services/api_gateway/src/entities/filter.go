package entities

import "time"

type Filter struct {
	CameraID      *int       `json:"camera_id"`
	ObjectID      *int       `json:"object_id"`
	Type          *string    `json:"type"`
	IP            *string    `json:"ip"`
	IsEnabled     *bool      `json:"is_enabled"`
	ConfidenceMin *float64   `json:"confidence_min"`
	ConfidenceMax *float64   `json:"confidence_max"`
	CreatedFrom   *time.Time `json:"created_from"`
	CreatedTo     *time.Time `json:"created_to"`
}
