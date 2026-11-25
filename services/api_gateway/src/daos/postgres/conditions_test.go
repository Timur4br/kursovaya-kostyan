package pg_dao

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"api_gateway/entities"
)

func TestEventFilterUsesParametersBeforePagination(t *testing.T) {
	cameraID, objectID := 2, 3
	eventType := "found' OR 1=1 --"
	confidence := 0.0
	from := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	filter := &entities.Filter{CameraID: &cameraID, ObjectID: &objectID, Type: &eventType, ConfidenceMin: &confidence, CreatedFrom: &from}
	args := []any{"existing"}
	query := addFilterByEvent("SELECT * FROM events WHERE photo_uuid = $1", &args, filter)
	query += " ORDER BY uuid"
	query = addPagination(query, &args, 10, 20)
	want := "SELECT * FROM events WHERE photo_uuid = $1 AND camera_id = $2 AND object_id = $3 AND type = $4 AND confidence >= $5 AND created_at >= $6 ORDER BY uuid LIMIT $7 OFFSET $8"
	if query != want {
		t.Fatalf("query = %q", query)
	}
	if strings.Contains(query, eventType) {
		t.Fatal("filter value interpolated into SQL")
	}
	if !reflect.DeepEqual(args, []any{"existing", cameraID, objectID, eventType, confidence, from, 10, 20}) {
		t.Fatalf("args = %#v", args)
	}
}

func TestCameraFalseAndObjectFilter(t *testing.T) {
	enabled := false
	args := []any{}
	query := addFilterByCamera("SELECT COUNT(*) FROM camera", &args, &entities.Filter{IsEnabled: &enabled})
	if query != "SELECT COUNT(*) FROM camera WHERE is_enabled = $1" || !reflect.DeepEqual(args, []any{false}) {
		t.Fatalf("query = %q, args = %#v", query, args)
	}
	objectType := "bag"
	args = nil
	query = addFilterByObject("SELECT COUNT(*) FROM objects", &args, &entities.Filter{Type: &objectType})
	if query != "SELECT COUNT(*) FROM objects WHERE type = $1" || !reflect.DeepEqual(args, []any{"bag"}) {
		t.Fatalf("query = %q, args = %#v", query, args)
	}
	for _, addFilter := range []func(string, *[]any, *entities.Filter) string{addFilterByCamera, addFilterByObject, addFilterByEvent} {
		args = nil
		if got := addFilter("SELECT 1", &args, nil); got != "SELECT 1" || len(args) != 0 {
			t.Fatal("absent filter changed query")
		}
	}
}
