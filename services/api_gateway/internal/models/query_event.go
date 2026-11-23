package models

type QueryEvent struct {
	EventID int64 `json:"event_id"`
	QueryID int64 `json:"query_id"`
}
