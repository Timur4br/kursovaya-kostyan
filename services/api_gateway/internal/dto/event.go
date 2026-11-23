package dto

import "encoding/json"

type EventInput struct {
	CameraID    Value[int64]   `json:"camera_id"`
	Confidence  Value[float64] `json:"confidence"`
	Description Value[string]  `json:"description"`
	EventTypeID Value[int64]   `json:"event_type_id"`
	LocationID  Value[int64]   `json:"location_id"`
	ObjectID    Value[int64]   `json:"object_id"`
	OccurredAt  Value[string]  `json:"occurred_at"`
	PersonID    Value[int64]   `json:"person_id"`
}

func (input EventInput) Fields() map[string]json.RawMessage {
	values := map[string]json.RawMessage{}
	put(values, "camera_id", input.CameraID)
	put(values, "confidence", input.Confidence)
	put(values, "description", input.Description)
	put(values, "event_type_id", input.EventTypeID)
	put(values, "location_id", input.LocationID)
	put(values, "object_id", input.ObjectID)
	put(values, "occurred_at", input.OccurredAt)
	put(values, "person_id", input.PersonID)
	return values
}
