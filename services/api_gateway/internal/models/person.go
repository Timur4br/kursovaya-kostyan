package models

import (
	"encoding/json"
	"time"
)

type Person struct {
	FirstSeen      *time.Time      `json:"first_seen"`
	ID             int64           `json:"id"`
	LastSeen       *time.Time      `json:"last_seen"`
	UserID         *int64          `json:"user_id"`
	VisualFeatures json.RawMessage `json:"visual_features"`
}
