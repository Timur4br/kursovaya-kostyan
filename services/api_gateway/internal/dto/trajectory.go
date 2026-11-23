package dto

import "encoding/json"

type TrajectoryInput struct {
	Direction       Value[string] `json:"direction"`
	EndLocationID   Value[int64]  `json:"end_location_id"`
	EndTime         Value[string] `json:"end_time"`
	PersonID        Value[int64]  `json:"person_id"`
	StartLocationID Value[int64]  `json:"start_location_id"`
	StartTime       Value[string] `json:"start_time"`
	Status          Value[string] `json:"status"`
}

func (input TrajectoryInput) Fields() map[string]json.RawMessage {
	values := map[string]json.RawMessage{}
	put(values, "direction", input.Direction)
	put(values, "end_location_id", input.EndLocationID)
	put(values, "end_time", input.EndTime)
	put(values, "person_id", input.PersonID)
	put(values, "start_location_id", input.StartLocationID)
	put(values, "start_time", input.StartTime)
	put(values, "status", input.Status)
	return values
}
