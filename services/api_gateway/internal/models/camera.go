package models

type Camera struct {
	FPS        *float64 `json:"fps"`
	ID         int64    `json:"id"`
	LocationID int64    `json:"location_id"`
	Name       string   `json:"name"`
	Resolution *string  `json:"resolution"`
	Status     string   `json:"status"`
	StreamURL  *string  `json:"stream_url"`
}
