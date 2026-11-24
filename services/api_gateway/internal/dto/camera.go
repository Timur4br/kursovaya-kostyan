package dto

import "encoding/json"

type CameraInput struct {
	FPS        Value[float64] `json:"fps"`
	LocationID Value[int64]   `json:"location_id"`
	Name       Value[string]  `json:"name"`
	Resolution Value[string]  `json:"resolution"`
	Status     Value[string]  `json:"status"`
	StreamURL  Value[string]  `json:"stream_url"`
}

func (input CameraInput) Fields() map[string]json.RawMessage {
	values := map[string]json.RawMessage{}
	put(values, "fps", input.FPS)
	put(values, "location_id", input.LocationID)
	put(values, "name", input.Name)
	put(values, "resolution", input.Resolution)
	put(values, "status", input.Status)
	put(values, "stream_url", input.StreamURL)
	return values
}
