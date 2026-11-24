package dto

import "encoding/json"

type QueryHistoryInput struct {
	CompletedAt Value[string] `json:"completed_at"`
	Response    Value[string] `json:"response"`
	Status      Value[string] `json:"status"`
	Text        Value[string] `json:"text"`
	UserID      Value[int64]  `json:"user_id"`
}

func (input QueryHistoryInput) Fields() map[string]json.RawMessage {
	values := map[string]json.RawMessage{}
	put(values, "completed_at", input.CompletedAt)
	put(values, "response", input.Response)
	put(values, "status", input.Status)
	put(values, "text", input.Text)
	put(values, "user_id", input.UserID)
	return values
}
