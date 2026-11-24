package dto

import "encoding/json"

type QueryEventInput struct {
	EventID Value[int64] `json:"event_id"`
	QueryID Value[int64] `json:"query_id"`
}

func (input QueryEventInput) Fields() map[string]json.RawMessage {
	values := map[string]json.RawMessage{}
	put(values, "event_id", input.EventID)
	put(values, "query_id", input.QueryID)
	return values
}
