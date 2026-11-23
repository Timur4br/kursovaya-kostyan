package app

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/config"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

func realAPI(t *testing.T) (*httptest.Server, *sql.DB) {
	t.Helper()
	if os.Getenv("RUN_DB_TESTS") != "1" {
		t.Skip("run through Compose tests profile for PostgreSQL CRUD checks")
	}
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	admin, err := database.Open(context.Background(), cfg.DatabaseURL)
	if err != nil {
		t.Fatal(err)
	}
	schema := fmt.Sprintf("api_test_%d", time.Now().UnixNano())
	quoted := pgx.Identifier{schema}.Sanitize()
	if _, err := admin.Exec("CREATE SCHEMA " + quoted); err != nil {
		admin.Close()
		t.Fatal(err)
	}
	conn, err := pgx.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		t.Fatal(err)
	}
	conn.RuntimeParams["search_path"] = schema
	db := stdlib.OpenDB(*conn)
	if err := database.Migrate(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(Handler(db, func(ctx context.Context) error { return database.Ready(ctx, db) }, slog.New(slog.NewTextHandler(io.Discard, nil))))
	t.Cleanup(func() {
		server.Close()
		db.Close()
		if _, err := admin.Exec("DROP SCHEMA " + quoted + " CASCADE"); err != nil {
			t.Error(err)
		}
		admin.Close()
	})
	return server, db
}

func TestCRUDWithPostgreSQL(t *testing.T) {
	server, db := realAPI(t)
	request := func(method, path string, body any, want int) map[string]json.RawMessage {
		t.Helper()
		var reader io.Reader
		if body != nil {
			encoded, err := json.Marshal(body)
			if err != nil {
				t.Fatal(err)
			}
			reader = bytes.NewReader(encoded)
		}
		req, err := http.NewRequest(method, server.URL+path, reader)
		if err != nil {
			t.Fatal(err)
		}
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		resp, err := server.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != want {
			t.Fatalf("%s %s: want %d, got %d: %s", method, path, want, resp.StatusCode, data)
		}
		if want == 204 {
			if len(data) != 0 {
				t.Fatal("204 has body")
			}
			return nil
		}
		var result map[string]json.RawMessage
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatal(err)
		}
		return result
	}
	type record struct {
		path  string
		id    int64
		patch map[string]any
	}
	records := []record{}
	create := func(path string, body, patch map[string]any) int64 {
		t.Helper()
		result := request("POST", "/"+path, body, 201)
		var row map[string]json.RawMessage
		if err := json.Unmarshal(result["data"], &row); err != nil {
			t.Fatal(err)
		}
		id, err := strconv.ParseInt(string(row["id"]), 10, 64)
		if err != nil {
			t.Fatal(err)
		}
		replacement := map[string]any{}
		for k, v := range body {
			replacement[k] = v
		}
		for k, v := range patch {
			replacement[k] = v
		}
		records = append(records, record{path, id, replacement})
		return id
	}
	u := create("user", map[string]any{"username": "owner", "name": "Owner"}, map[string]any{"name": "Updated Owner"})
	loc := create("location", map[string]any{"name": "Office", "type": "room", "coordinates": map[string]any{"x": 1}}, map[string]any{"description": "Updated office"})
	otherLoc := create("location", map[string]any{"name": "Kitchen", "type": "room"}, map[string]any{"description": "Other room"})
	cam := create("camera", map[string]any{"location_id": loc, "name": "Camera", "fps": 25}, map[string]any{"status": "online"})
	person := create("person", map[string]any{"user_id": u, "visual_features": map[string]any{"label": "test"}}, map[string]any{"visual_features": map[string]any{"label": "updated"}})
	objectType := create("object_type", map[string]any{"name": "Backpack"}, map[string]any{"description": "Updated type"})
	object := create("object", map[string]any{"object_type_id": objectType, "name": "My backpack"}, map[string]any{"status": "visible"})
	eventType := create("event_type", map[string]any{"name": "placed"}, map[string]any{"description": "Placed on a desk"})
	event := create("event", map[string]any{"event_type_id": eventType, "camera_id": cam, "location_id": loc, "person_id": person, "object_id": object, "occurred_at": "2026-11-22T10:00:00Z", "confidence": 0.95}, map[string]any{"confidence": 0.9})
	video := create("video_segment", map[string]any{"camera_id": cam, "event_id": event, "start_time": "2026-11-22T10:00:00Z", "end_time": "2026-11-22T10:00:12Z", "storage_reference": "test.mp4"}, map[string]any{"preview": "preview.jpg"})
	trajectory := create("trajectory", map[string]any{"person_id": person, "start_time": "2026-11-22T10:00:00Z", "start_location_id": loc}, map[string]any{"direction": "east"})
	create("trajectory_point", map[string]any{"trajectory_id": trajectory, "video_segment_id": video, "observed_at": "2026-11-22T10:00:02Z", "x": 1, "y": 2}, map[string]any{"z": 3})
	query := create("query_history", map[string]any{"user_id": u, "text": "Where is my backpack?"}, map[string]any{"status": "completed", "response": "In the office"})
	link := map[string]any{"query_id": query, "event_id": event}
	request("POST", "/query_event", link, 201)
	linkPath := fmt.Sprintf("/query_event/%d/%d", query, event)
	request("GET", linkPath, nil, 200)
	request("POST", "/query_event", link, 409)
	for _, row := range records {
		t.Run(row.path+"/"+strconv.FormatInt(row.id, 10), func(t *testing.T) {
			path := "/" + row.path + "/" + strconv.FormatInt(row.id, 10)
			request("GET", path, nil, 200)
			updated := request("PUT", path, row.patch, 200)
			var data map[string]any
			if err := json.Unmarshal(updated["data"], &data); err != nil {
				t.Fatal(err)
			}
			for key, value := range row.patch {
				actual, _ := json.Marshal(data[key])
				expected, _ := json.Marshal(value)
				if string(actual) != string(expected) {
					t.Fatalf("patch %s did not persist: %s vs %s", key, actual, expected)
				}
			}
		})
	}
	// Credential storage exists, but neither reads nor CRUD writes expose it.
	if _, err := db.Exec("UPDATE users SET password_hash='DO_NOT_EXPOSE' WHERE id=$1", u); err != nil {
		t.Fatal(err)
	}
	user := request("GET", fmt.Sprintf("/user/%d", u), nil, 200)
	if strings.Contains(string(user["data"]), "password_hash") || strings.Contains(string(user["data"]), "DO_NOT_EXPOSE") {
		t.Fatal("credential leaked")
	}
	request("PUT", fmt.Sprintf("/user/%d", u), map[string]any{"password_hash": "changed"}, 422)
	request("POST", "/user", map[string]any{"username": "owner", "name": "Duplicate"}, 409)
	request("POST", "/camera", map[string]any{"name": "Bad", "location_id": 999999}, 409)
	request("PUT", fmt.Sprintf("/camera/%d", cam), map[string]any{"location_id": otherLoc, "name": "Moved camera"}, 409)
	request("PUT", fmt.Sprintf("/video_segment/%d", video), map[string]any{"camera_id": cam, "event_id": event, "start_time": "2026-11-22T10:00:00Z", "end_time": "2026-11-22T09:59:00Z", "storage_reference": "test.mp4"}, 422)
	request("DELETE", fmt.Sprintf("/camera/%d", cam), nil, 409)
	request("DELETE", fmt.Sprintf("/user/%d", u), nil, 409)
	for _, method := range []string{"GET", "DELETE", "PUT"} {
		var body any
		if method == "PUT" {
			body = map[string]any{"name": "Missing", "type": "room"}
		}
		request(method, "/location/999999", body, 404)
	}
	// PUT resets omitted optional fields and keeps immutable fields.
	reset := request("PUT", fmt.Sprintf("/camera/%d", cam), map[string]any{"location_id": loc, "name": "Reset camera"}, 200)
	var resetData map[string]json.RawMessage
	if err := json.Unmarshal(reset["data"], &resetData); err != nil {
		t.Fatal(err)
	}
	if string(resetData["fps"]) != "null" || string(resetData["status"]) != `"offline"` {
		t.Fatal("omitted optional fields were not reset", string(reset["data"]))
	}
	request("PUT", fmt.Sprintf("/person/%d", person), map[string]any{}, 200)
	cleared := request("GET", fmt.Sprintf("/person/%d", person), nil, 200)
	var clearedData map[string]json.RawMessage
	if err := json.Unmarshal(cleared["data"], &clearedData); err != nil {
		t.Fatal(err)
	}
	if string(clearedData["user_id"]) != "null" || string(clearedData["visual_features"]) != "null" {
		t.Fatal("PUT did not clear nullable fields")
	}
	// Missing required values make a full replacement invalid.
	request("PUT", fmt.Sprintf("/user/%d", u), map[string]any{"name": "Partial update"}, 422)
	// A query/event pair also has all four operations; PUT can replace its composite key.
	replacementQuery := create("query_history", map[string]any{"user_id": u, "text": "Another question"}, map[string]any{"response": "answer"})
	request("PUT", linkPath, map[string]any{"query_id": replacementQuery, "event_id": event}, 200)
	request("GET", linkPath, nil, 404)
	linkPath = fmt.Sprintf("/query_event/%d/%d", replacementQuery, event)
	request("GET", linkPath, nil, 200)
	// Values containing SQL punctuation are ordinary bound parameters.
	attack := create("user", map[string]any{"username": "x'; DROP TABLE cameras; --", "name": "Injection test"}, map[string]any{"name": "Safe"})
	request("DELETE", fmt.Sprintf("/user/%d", attack), nil, 204)
	records = records[:len(records)-1]
	request("GET", fmt.Sprintf("/camera/%d", cam), nil, 200)
	request("DELETE", linkPath, nil, 204)
	request("GET", linkPath, nil, 404)
	for i := len(records) - 1; i >= 0; i-- {
		row := records[i]
		path := fmt.Sprintf("/%s/%d", row.path, row.id)
		request("DELETE", path, nil, 204)
		request("GET", path, nil, 404)
	}
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	value, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return value
}
