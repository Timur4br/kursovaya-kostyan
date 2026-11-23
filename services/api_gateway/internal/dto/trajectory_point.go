package dto

import "encoding/json"

type TrajectoryPointInput struct {
	ObservedAt     Value[string]  `json:"observed_at"`
	TrajectoryID   Value[int64]   `json:"trajectory_id"`
	VideoSegmentID Value[int64]   `json:"video_segment_id"`
	X              Value[float64] `json:"x"`
	Y              Value[float64] `json:"y"`
	Z              Value[float64] `json:"z"`
}

func (input TrajectoryPointInput) Fields() map[string]json.RawMessage {
	values := map[string]json.RawMessage{}
	put(values, "observed_at", input.ObservedAt)
	put(values, "trajectory_id", input.TrajectoryID)
	put(values, "video_segment_id", input.VideoSegmentID)
	put(values, "x", input.X)
	put(values, "y", input.Y)
	put(values, "z", input.Z)
	return values
}
