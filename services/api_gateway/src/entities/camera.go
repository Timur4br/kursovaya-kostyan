package entities

import "time"

type Camera struct {
	ID        int       `json:"id"`
	Comment   string    `json:"comment"`
	IP        string    `json:"ip"`
	Port      int       `json:"port"`
	StreamURL string    `json:"stream_url"`
	FPS       int       `json:"fps"`
	IsEnabled bool      `json:"is_enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CameraUpdate distinguishes omitted fields from explicit empty strings or false.
type CameraUpdate struct {
	Comment   *string `json:"comment"`
	IP        *string `json:"ip" binding:"omitempty,min=1"`
	Port      *int    `json:"port" binding:"omitempty,gte=1,lte=65535"`
	StreamURL *string `json:"stream_url"`
	IsEnabled *bool   `json:"is_enabled"`
}
