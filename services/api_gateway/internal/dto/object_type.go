package dto

import "encoding/json"

type ObjectTypeInput struct {
	Description Value[string] `json:"description"`
	Name        Value[string] `json:"name"`
}

func (input ObjectTypeInput) Fields() map[string]json.RawMessage {
	values := map[string]json.RawMessage{}
	put(values, "description", input.Description)
	put(values, "name", input.Name)
	return values
}
