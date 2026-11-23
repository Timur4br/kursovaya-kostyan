package models

import (
	"encoding/json"
)

type Location struct {
	Coordinates json.RawMessage `json:"coordinates"`
	Description string          `json:"description"`
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Type        string          `json:"type"`
}
