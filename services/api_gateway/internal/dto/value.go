package dto

import "encoding/json"

// Value preserves absent vs explicit null while decoding typed request fields.
type Value[T any] struct {
	Value   T
	Present bool
	Null    bool
	raw     json.RawMessage
}

func (v *Value[T]) UnmarshalJSON(data []byte) error {
	v.Present = true
	v.raw = append(v.raw[:0], data...)
	if string(data) == "null" {
		v.Null = true
		return nil
	}
	v.Null = false
	return json.Unmarshal(data, &v.Value)
}
func put[T any](values map[string]json.RawMessage, name string, v Value[T]) {
	if v.Present {
		values[name] = v.raw
	}
}
