// Package resources defines the API's explicit allowlist of tables and columns.
// SQL identifiers are never taken from request bodies or query strings.
package resources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type Kind string

const (
	Text      Kind = "string"
	ID        Kind = "id"
	Number    Kind = "number"
	Timestamp Kind = "timestamp"
	Object    Kind = "object"
)

type Field struct {
	Name                                 string
	Kind                                 Kind
	Required, Nullable, ReadOnly, Filter bool
	MaxLength                            int
	Enum                                 []string
	Positive, Confidence                 bool
}

type Resource struct {
	Path, Table string
	Keys        []string
	Fields      []Field
}

func id() Field { return Field{Name: "id", Kind: ID, ReadOnly: true} }
func ref(name string, nullable bool) Field {
	return Field{Name: name, Kind: ID, Required: !nullable, Nullable: nullable, Filter: true}
}
func text(name string, required bool, max int) Field {
	return Field{Name: name, Kind: Text, Required: required, MaxLength: max}
}
func nullableText(name string) Field {
	return Field{Name: name, Kind: Text, Nullable: true, MaxLength: 4096}
}
func timestamp(name string, required bool) Field {
	return Field{Name: name, Kind: Timestamp, Required: required, Nullable: !required}
}
func created() Field { return Field{Name: "created_at", Kind: Timestamp, ReadOnly: true} }
func status(values ...string) Field {
	return Field{Name: "status", Kind: Text, Filter: true, Enum: values}
}
func object(name string) Field { return Field{Name: name, Kind: Object, Nullable: true} }
func entity(path, table string, fields ...Field) Resource {
	return Resource{Path: path, Table: table, Keys: []string{"id"}, Fields: append([]Field{id()}, fields...)}
}

var All = []Resource{
	entity("user", "users", text("username", true, 100), text("name", true, 200), created()),
	entity("location", "locations", text("name", true, 200), text("type", true, 100), text("description", false, 16000), object("coordinates")),
	entity("camera", "cameras", ref("location_id", false), text("name", true, 200), nullableText("stream_url"), nullableText("resolution"), Field{Name: "fps", Kind: Number, Nullable: true, Positive: true}, status("offline", "online", "disabled")),
	entity("person", "persons", ref("user_id", true), object("visual_features"), timestamp("first_seen", false), timestamp("last_seen", false)),
	entity("object_type", "object_types", text("name", true, 100), text("description", false, 16000)),
	entity("object", "objects", ref("object_type_id", false), text("name", true, 200), status("unknown", "visible", "lost", "found"), nullableText("image_path"), timestamp("first_seen", false), timestamp("last_seen", false)),
	entity("event_type", "event_types", text("name", true, 100), text("description", false, 16000)),
	entity("event", "events", ref("event_type_id", false), ref("camera_id", false), ref("location_id", false), ref("person_id", true), ref("object_id", true), timestamp("occurred_at", true), Field{Name: "confidence", Kind: Number, Nullable: true, Confidence: true}, text("description", false, 16000)),
	entity("video_segment", "video_segments", ref("camera_id", false), ref("event_id", true), timestamp("start_time", true), timestamp("end_time", true), text("storage_reference", true, 4096), nullableText("preview"), Field{Name: "duration_seconds", Kind: Number, ReadOnly: true}),
	entity("trajectory", "trajectories", ref("person_id", false), timestamp("start_time", true), timestamp("end_time", false), ref("start_location_id", true), ref("end_location_id", true), status("active", "completed", "lost"), nullableText("direction")),
	entity("trajectory_point", "trajectory_points", ref("trajectory_id", false), ref("video_segment_id", true), timestamp("observed_at", true), Field{Name: "x", Kind: Number, Required: true}, Field{Name: "y", Kind: Number, Required: true}, Field{Name: "z", Kind: Number, Nullable: true}),
	entity("query_history", "query_history", ref("user_id", false), text("text", true, 4000), Field{Name: "response", Kind: Text, Nullable: true, MaxLength: 64000}, status("pending", "processing", "completed", "failed"), created(), timestamp("completed_at", false)),
	{Path: "query_event", Table: "query_events", Keys: []string{"query_id", "event_id"}, Fields: []Field{ref("query_id", false), ref("event_id", false)}},
}

func (r Resource) Field(name string) (Field, bool) {
	for _, f := range r.Fields {
		if f.Name == name {
			return f, true
		}
	}
	return Field{}, false
}
func (r Resource) CanUpdate() bool { return len(r.Keys) == 1 }
func (r Resource) ItemPath() string {
	path := "/" + r.Path
	for _, key := range r.Keys {
		path += "/{" + key + "}"
	}
	return path
}

type ValidationError struct{ Field, Message string }

func (e *ValidationError) Error() string  { return e.Field + ": " + e.Message }
func invalid(field, message string) error { return &ValidationError{Field: field, Message: message} }

func (r Resource) ValidateBody(raw map[string]json.RawMessage, patch bool) (map[string]any, error) {

	result := make(map[string]any, len(raw))
	for name, value := range raw {
		field, ok := r.Field(name)
		if !ok {
			return nil, invalid(name, "неизвестное поле")
		}
		if field.ReadOnly {
			return nil, invalid(name, "поле только для чтения")
		}
		parsed, err := field.Parse(value)
		if err != nil {
			return nil, err
		}
		result[name] = parsed
	}
	if !patch {
		for _, field := range r.Fields {
			if field.Required {
				if _, ok := raw[field.Name]; !ok {
					return nil, invalid(field.Name, "обязательное поле")
				}
			}
		}
	}
	return result, nil
}

func (f Field) Parse(raw json.RawMessage) (any, error) {
	if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		if f.Nullable {
			return nil, nil
		}
		return nil, invalid(f.Name, "не может быть null")
	}
	switch f.Kind {
	case Text:
		var value string
		if err := json.Unmarshal(raw, &value); err != nil {
			return nil, invalid(f.Name, "ожидается строка")
		}
		if f.MaxLength > 0 && utf8.RuneCountInString(value) > f.MaxLength {
			return nil, invalid(f.Name, fmt.Sprintf("максимум %d символов", f.MaxLength))
		}
		if f.Required && strings.TrimSpace(value) == "" {
			return nil, invalid(f.Name, "не может быть пустым")
		}
		if len(f.Enum) > 0 {
			for _, allowed := range f.Enum {
				if value == allowed {
					return value, nil
				}
			}
			return nil, invalid(f.Name, "недопустимое значение")
		}
		return value, nil
	case ID:
		value, err := strconv.ParseInt(string(raw), 10, 64)
		if err != nil || value <= 0 {
			return nil, invalid(f.Name, "ожидается положительное целое число")
		}
		return value, nil
	case Number:
		value, err := strconv.ParseFloat(string(raw), 64)
		if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
			return nil, invalid(f.Name, "ожидается конечное число")
		}
		if f.Positive && value <= 0 {
			return nil, invalid(f.Name, "должно быть больше нуля")
		}
		if f.Confidence && (value < 0 || value > 1) {
			return nil, invalid(f.Name, "ожидается число от 0 до 1")
		}
		return value, nil
	case Timestamp:
		var value string
		if err := json.Unmarshal(raw, &value); err != nil {
			return nil, invalid(f.Name, "ожидается дата RFC3339 с часовым поясом")
		}
		parsed, err := time.Parse(time.RFC3339Nano, value)
		if err != nil {
			return nil, invalid(f.Name, "ожидается дата RFC3339 с часовым поясом")
		}
		return parsed.UTC(), nil
	case Object:
		if len(bytes.TrimSpace(raw)) == 0 || bytes.TrimSpace(raw)[0] != '{' || !json.Valid(raw) {
			return nil, invalid(f.Name, "ожидается JSON-объект")
		}
		return string(raw), nil
	}
	return nil, invalid(f.Name, "неподдерживаемый тип")
}

// Must returns a compiled-in resource definition, never a request-supplied table.
func Must(path string) Resource {
	for _, r := range All {
		if r.Path == path {
			return r
		}
	}
	panic("unknown compiled resource: " + path)
}
