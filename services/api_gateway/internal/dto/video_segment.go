package dto

import "encoding/json"

type VideoSegmentInput struct {
	CameraID         Value[int64]  `json:"camera_id"`
	EndTime          Value[string] `json:"end_time"`
	EventID          Value[int64]  `json:"event_id"`
	Preview          Value[string] `json:"preview"`
	StartTime        Value[string] `json:"start_time"`
	StorageReference Value[string] `json:"storage_reference"`
}

func (input VideoSegmentInput) Fields() map[string]json.RawMessage {
	values := map[string]json.RawMessage{}
	put(values, "camera_id", input.CameraID)
	put(values, "end_time", input.EndTime)
	put(values, "event_id", input.EventID)
	put(values, "preview", input.Preview)
	put(values, "start_time", input.StartTime)
	put(values, "storage_reference", input.StorageReference)
	return values
}
