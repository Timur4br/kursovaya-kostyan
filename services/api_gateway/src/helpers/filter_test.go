package helpers

import (
	"net/url"
	"testing"
	"time"
)

func TestExtractFilterFromQueryParams(t *testing.T) {
	filter, err := ExtractFilterFromQueryParams(url.Values{"limit": {"5"}})
	if err != nil || filter != nil {
		t.Fatalf("absent filter: %v, %v", filter, err)
	}
	params := url.Values{
		"camera_id": {"3"}, "object_id": {"4"}, "type": {"bag' OR 1=1 --"},
		"ip": {"127.0.0.1"}, "is_enabled": {"false"},
		"confidence_min": {"0"}, "confidence_max": {"0.9"},
		"created_from": {"2026-11-01T12:00:00+03:00"}, "created_to": {"2026-11-02"},
	}
	filter, err = ExtractFilterFromQueryParams(params)
	if err != nil {
		t.Fatal(err)
	}
	if filter == nil || filter.CameraID == nil || *filter.CameraID != 3 || filter.ObjectID == nil || *filter.ObjectID != 4 ||
		filter.Type == nil || *filter.Type != params.Get("type") || filter.IP == nil || *filter.IP != "127.0.0.1" ||
		filter.IsEnabled == nil || *filter.IsEnabled || filter.ConfidenceMin == nil || *filter.ConfidenceMin != 0 ||
		filter.ConfidenceMax == nil || *filter.ConfidenceMax != 0.9 {
		t.Fatalf("filter did not preserve supplied values: %+v", filter)
	}
	if filter.CreatedFrom == nil || !filter.CreatedFrom.Equal(time.Date(2026, 11, 1, 9, 0, 0, 0, time.UTC)) ||
		filter.CreatedTo == nil || !filter.CreatedTo.Equal(time.Date(2026, 11, 2, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("creation time bounds were not parsed correctly")
	}
}

func TestExtractFilterRejectsInvalidParameters(t *testing.T) {
	for _, query := range []string{
		"camera_id=bad", "object_id=1.5", "is_enabled=maybe", "confidence_min=bad",
		"confidence_min=NaN", "confidence_max=Inf", "confidence_min=0.9&confidence_max=0.2",
		"created_from=bad", "created_to=2026-02-30", "created_from=2026-11-02&created_to=2026-11-01",
	} {
		t.Run(query, func(t *testing.T) {
			params, err := url.ParseQuery(query)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := ExtractFilterFromQueryParams(params); err == nil {
				t.Fatal("expected invalid filter error")
			}
		})
	}
}
