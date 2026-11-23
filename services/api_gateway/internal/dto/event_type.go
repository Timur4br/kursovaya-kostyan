package dto

import "encoding/json"

type EventTypeInput struct {
	Description Value[string] `json:"description"`
	Name        Value[string] `json:"name"`
}

func (input EventTypeInput) Fields() map[string]json.RawMessage {
	values := map[string]json.RawMessage{}
	put(values, "description", input.Description)
	put(values, "name", input.Name)
	return values
}
