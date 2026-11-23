package dto

import "encoding/json"

type UserInput struct {
	Name     Value[string] `json:"name"`
	Username Value[string] `json:"username"`
}

func (input UserInput) Fields() map[string]json.RawMessage {
	values := map[string]json.RawMessage{}
	put(values, "name", input.Name)
	put(values, "username", input.Username)
	return values
}
