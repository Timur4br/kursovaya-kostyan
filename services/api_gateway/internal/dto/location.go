package dto

import "encoding/json"

type LocationInput struct {
	Coordinates Value[json.RawMessage] `json:"coordinates"`
	Description Value[string]          `json:"description"`
	Name        Value[string]          `json:"name"`
	Type        Value[string]          `json:"type"`
}

func (input LocationInput) Fields() map[string]json.RawMessage {
	values := map[string]json.RawMessage{}
	put(values, "coordinates", input.Coordinates)
	put(values, "description", input.Description)
	put(values, "name", input.Name)
	put(values, "type", input.Type)
	return values
}
