package dto

import "encoding/json"

type PersonInput struct {
	FirstSeen      Value[string]          `json:"first_seen"`
	LastSeen       Value[string]          `json:"last_seen"`
	UserID         Value[int64]           `json:"user_id"`
	VisualFeatures Value[json.RawMessage] `json:"visual_features"`
}

func (input PersonInput) Fields() map[string]json.RawMessage {
	values := map[string]json.RawMessage{}
	put(values, "first_seen", input.FirstSeen)
	put(values, "last_seen", input.LastSeen)
	put(values, "user_id", input.UserID)
	put(values, "visual_features", input.VisualFeatures)
	return values
}
