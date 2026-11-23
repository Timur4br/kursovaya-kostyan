package dto

import "encoding/json"

type ObjectInput struct {
	FirstSeen    Value[string] `json:"first_seen"`
	ImagePath    Value[string] `json:"image_path"`
	LastSeen     Value[string] `json:"last_seen"`
	Name         Value[string] `json:"name"`
	ObjectTypeID Value[int64]  `json:"object_type_id"`
	Status       Value[string] `json:"status"`
}

func (input ObjectInput) Fields() map[string]json.RawMessage {
	values := map[string]json.RawMessage{}
	put(values, "first_seen", input.FirstSeen)
	put(values, "image_path", input.ImagePath)
	put(values, "last_seen", input.LastSeen)
	put(values, "name", input.Name)
	put(values, "object_type_id", input.ObjectTypeID)
	put(values, "status", input.Status)
	return values
}
