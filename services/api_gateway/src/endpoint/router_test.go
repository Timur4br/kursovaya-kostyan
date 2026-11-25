package endpoint

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"api_gateway/entities"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func request(t *testing.T, router http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var data []byte
	var err error
	if body != nil {
		data, err = json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	return response
}

func expectStatus(t *testing.T, response *httptest.ResponseRecorder, status int) {
	t.Helper()
	if response.Code != status {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, status, response.Body.String())
	}
}

func decode[T any](t *testing.T, response *httptest.ResponseRecorder) T {
	t.Helper()
	var result T
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	return result
}

func TestEntityInvalidRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(nil).Router
	for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodDelete} {
		for _, path := range []string{"/api/object/invalid", "/api/event/invalid"} {
			t.Run(method+path, func(t *testing.T) {
				expectStatus(t, request(t, router, method, path, nil), http.StatusBadRequest)
			})
		}
	}
	event := entities.Event{
		CameraID: 1, ObjectID: 1, Type: "found", Confidence: 0,
		Embedding: make([]float32, 511), PhotoUUID: "bbf2ad34-030d-48c9-835c-c106bb6b516b",
	}
	expectStatus(t, request(t, router, http.MethodPost, "/api/object", map[string]any{}), http.StatusBadRequest)
	expectStatus(t, request(t, router, http.MethodPost, "/api/event", event), http.StatusBadRequest)
	event.Embedding = make([]float32, 512)
	event.PhotoUUID = "invalid"
	expectStatus(t, request(t, router, http.MethodPost, "/api/event", event), http.StatusBadRequest)
	expectStatus(t, request(t, router, http.MethodPut, "/api/event/bbf2ad34-030d-48c9-835c-c106bb6b516b", event), http.StatusBadRequest)
}

// TEST_DATABASE_URL must point to PostgreSQL with the vector extension installed.
// The schema comes from init-db.sql; tables are temporary and disappear on close.
func TestEntityPostgresCRUD(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_URL to run PostgreSQL/pgvector integration tests")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	if err := db.Ping(); err != nil {
		t.Fatal("cannot connect to test database")
	}
	schema, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "scripts", "init-db.sql"))
	if err != nil {
		t.Fatal(err)
	}
	temporarySchema := strings.ReplaceAll(string(schema), "CREATE EXTENSION IF NOT EXISTS vector;", "")
	temporarySchema = strings.ReplaceAll(temporarySchema, "CREATE TABLE ", "CREATE TEMP TABLE ")
	if _, err := db.Exec(temporarySchema); err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)
	router := NewRouter(db).Router
	cameraResponse := request(t, router, http.MethodPost, "/api/camera", entities.Camera{IP: "127.0.0.1", Port: 80, IsEnabled: true})
	expectStatus(t, cameraResponse, http.StatusCreated)
	if id := decode[int](t, cameraResponse); id != 1 {
		t.Fatalf("camera ID = %d", id)
	}
	cameraResponse = request(t, router, http.MethodGet, "/api/camera/1", nil)
	expectStatus(t, cameraResponse, http.StatusOK)
	originalCamera := decode[entities.Camera](t, cameraResponse)
	cameraResponse = request(t, router, http.MethodPut, "/api/camera/1", map[string]any{
		"id": 999, "comment": "Edited camera", "ip": "127.0.0.10", "port": 8554,
		"stream_url": "rtsp://127.0.0.10/live", "is_enabled": false,
	})
	expectStatus(t, cameraResponse, http.StatusOK)
	cameraResponse = request(t, router, http.MethodGet, "/api/camera/1", nil)
	expectStatus(t, cameraResponse, http.StatusOK)
	editedCamera := decode[entities.Camera](t, cameraResponse)
	if editedCamera.ID != 1 || editedCamera.Comment != "Edited camera" || editedCamera.IP != "127.0.0.10" ||
		editedCamera.Port != 8554 || editedCamera.StreamURL != "rtsp://127.0.0.10/live" || editedCamera.IsEnabled ||
		!editedCamera.CreatedAt.Equal(originalCamera.CreatedAt) {
		t.Fatalf("camera edit did not persist: %+v", editedCamera)
	}
	expectStatus(t, request(t, router, http.MethodPut, "/api/camera/1", map[string]any{"comment": "", "stream_url": ""}), http.StatusOK)
	cameraResponse = request(t, router, http.MethodGet, "/api/camera/1", nil)
	expectStatus(t, cameraResponse, http.StatusOK)
	partialCamera := decode[entities.Camera](t, cameraResponse)
	if partialCamera.Comment != "" || partialCamera.StreamURL != "" || partialCamera.IP != "127.0.0.10" ||
		partialCamera.Port != 8554 || partialCamera.IsEnabled {
		t.Fatalf("partial update lost fields: %+v", partialCamera)
	}
	expectStatus(t, request(t, router, http.MethodPut, "/api/camera/1", map[string]any{"ip": "127.0.0.1", "port": 80, "is_enabled": true}), http.StatusOK)
	for _, path := range []string{"/api/object", "/api/event"} {
		response := request(t, router, http.MethodGet, path+"/list?limit=10&offset=0", nil)
		expectStatus(t, response, http.StatusOK)
		if response.Body.String() != "[]" {
			t.Fatalf("empty list: %s", response.Body.String())
		}
	}

	response := request(t, router, http.MethodPost, "/api/object", entities.Object{Type: "bag"})
	expectStatus(t, response, http.StatusCreated)
	objectID := decode[int](t, response)
	if objectID != 1 {
		t.Fatalf("generated object ID = %d", objectID)
	}
	response = request(t, router, http.MethodGet, "/api/object/1", nil)
	expectStatus(t, response, http.StatusOK)
	object := decode[entities.Object](t, response)
	if object.Type != "bag" || object.CreatedAt.IsZero() {
		t.Fatalf("object = %+v", object)
	}
	response = request(t, router, http.MethodPut, "/api/object/1", entities.Object{ID: 999, Type: "backpack"})
	expectStatus(t, response, http.StatusOK)
	response = request(t, router, http.MethodGet, "/api/object/1", nil)
	expectStatus(t, response, http.StatusOK)
	updatedObject := decode[entities.Object](t, response)
	if updatedObject.ID != 1 || updatedObject.Type != "backpack" || !updatedObject.CreatedAt.Equal(object.CreatedAt) {
		t.Fatalf("updated object = %+v", updatedObject)
	}
	response = request(t, router, http.MethodGet, "/api/object/count", nil)
	expectStatus(t, response, http.StatusOK)
	if decode[int](t, response) != 1 {
		t.Fatal("object count")
	}
	response = request(t, router, http.MethodGet, "/api/object/list?limit=1&offset=0", nil)
	expectStatus(t, response, http.StatusOK)
	if objects := decode[[]entities.Object](t, response); len(objects) != 1 || objects[0].Type != "backpack" {
		t.Fatalf("objects = %+v", objects)
	}

	event := entities.Event{
		CameraID: 1, ObjectID: objectID, Type: "found", Confidence: 0.75,
		Embedding: make([]float32, 512), PhotoUUID: "bbf2ad34-030d-48c9-835c-c106bb6b516b",
	}
	event.Embedding[0] = 0.125
	event.Embedding[511] = -0.25
	response = request(t, router, http.MethodPost, "/api/event", event)
	expectStatus(t, response, http.StatusCreated)
	uuid := decode[string](t, response)
	response = request(t, router, http.MethodGet, "/api/event/"+uuid, nil)
	expectStatus(t, response, http.StatusOK)
	storedEvent := decode[entities.Event](t, response)
	if storedEvent.UUID != uuid || storedEvent.Confidence != event.Confidence || storedEvent.CreatedAt.IsZero() ||
		len(storedEvent.Embedding) != 512 || storedEvent.Embedding[0] != 0.125 || storedEvent.Embedding[511] != -0.25 {
		t.Fatalf("event did not round trip: %+v", storedEvent)
	}
	event.UUID = "00000000-0000-0000-0000-000000000000"
	event.Type = "updated"
	event.Confidence = 0
	event.Embedding[0] = 0
	response = request(t, router, http.MethodPut, "/api/event/"+uuid, event)
	expectStatus(t, response, http.StatusOK)
	response = request(t, router, http.MethodGet, "/api/event/"+uuid, nil)
	expectStatus(t, response, http.StatusOK)
	updatedEvent := decode[entities.Event](t, response)
	if updatedEvent.UUID != uuid || updatedEvent.Type != "updated" || updatedEvent.Confidence != 0 ||
		updatedEvent.Embedding[0] != 0 || !updatedEvent.CreatedAt.Equal(storedEvent.CreatedAt) {
		t.Fatalf("updated event = %+v", updatedEvent)
	}
	response = request(t, router, http.MethodGet, "/api/event/count", nil)
	expectStatus(t, response, http.StatusOK)
	if decode[int](t, response) != 1 {
		t.Fatal("event count")
	}
	response = request(t, router, http.MethodGet, "/api/event/list?limit=1&offset=0", nil)
	expectStatus(t, response, http.StatusOK)
	if events := decode[[]entities.Event](t, response); len(events) != 1 || events[0].UUID != uuid || events[0].Type != "updated" {
		t.Fatalf("events = %+v", events)
	}
	for _, path := range []string{"/api/object", "/api/event"} {
		response = request(t, router, http.MethodGet, path+"/list?limit=1&offset=1", nil)
		expectStatus(t, response, http.StatusOK)
		if response.Body.String() != "[]" {
			t.Fatalf("pagination: %s", response.Body.String())
		}
	}
	// Add contrasting rows so a missing or ignored filter changes the result.
	if _, err := db.Exec(`INSERT INTO camera (ip, is_enabled) VALUES ('127.0.0.2', FALSE)`); err != nil {
		t.Fatal(err)
	}
	response = request(t, router, http.MethodPost, "/api/object", entities.Object{Type: "keys"})
	expectStatus(t, response, http.StatusCreated)
	event.CameraID, event.ObjectID, event.Type, event.Confidence = 2, 2, "other", 0.9
	response = request(t, router, http.MethodPost, "/api/event", event)
	expectStatus(t, response, http.StatusCreated)
	otherUUID := decode[string](t, response)
	response = request(t, router, http.MethodGet, "/api/event/list?limit=1&offset=0", nil)
	expectStatus(t, response, http.StatusOK)
	if events := decode[[]entities.Event](t, response); len(events) != 1 || events[0].UUID != otherUUID {
		t.Fatal("newest event must come first")
	}
	response = request(t, router, http.MethodGet, "/api/event/list?limit=1&offset=1", nil)
	expectStatus(t, response, http.StatusOK)
	if events := decode[[]entities.Event](t, response); len(events) != 1 || events[0].UUID != uuid {
		t.Fatal("older event must be on the next page")
	}
	for _, tc := range []struct {
		path, query string
		count       int
	}{
		{"camera", "camera_id=1", 1},
		{"camera", "ip=127.0.0.2&is_enabled=false", 1},
		{"camera", "camera_id=1&is_enabled=false", 0},
		{"object", "object_id=2&type=keys", 1},
		{"object", "type=backpack", 1},
		{"object", "type=" + url.QueryEscape("backpack' OR 1=1 --"), 0},
		{"event", "camera_id=1&object_id=1&type=updated&confidence_min=0&confidence_max=0", 1},
		{"event", "camera_id=1&object_id=2", 0},
		{"event", "confidence_min=0.8", 1},
		{"event", "confidence_max=0", 1},
		{"event", "created_from=" + url.QueryEscape(storedEvent.CreatedAt.Format("2006-01-02T15:04:05.999999999Z07:00")) +
			"&created_to=" + url.QueryEscape(storedEvent.CreatedAt.Format("2006-01-02T15:04:05.999999999Z07:00")) + "&camera_id=1", 1},
		{"camera", "created_from=9999-01-01", 0},
		{"object", "created_from=9999-01-01", 0},
		{"event", "created_from=9999-01-01", 0},
	} {
		t.Run(tc.path+"?"+tc.query, func(t *testing.T) {
			response := request(t, router, http.MethodGet, "/api/"+tc.path+"/count?"+tc.query, nil)
			expectStatus(t, response, http.StatusOK)
			if count := decode[int](t, response); count != tc.count {
				t.Fatalf("count = %d, want %d", count, tc.count)
			}
			response = request(t, router, http.MethodGet, "/api/"+tc.path+"/list?limit=10&offset=0&"+tc.query, nil)
			expectStatus(t, response, http.StatusOK)
			if items := decode[[]json.RawMessage](t, response); len(items) != tc.count {
				t.Fatalf("list size = %d, want %d", len(items), tc.count)
			}
			response = request(t, router, http.MethodGet, "/api/"+tc.path+"/list?limit=1&offset=1&"+tc.query, nil)
			expectStatus(t, response, http.StatusOK)
			if items := decode[[]json.RawMessage](t, response); len(items) != 0 {
				t.Fatal("filter must precede pagination")
			}
		})
	}
	expectStatus(t, request(t, router, http.MethodDelete, "/api/event/"+otherUUID, nil), http.StatusOK)
	expectStatus(t, request(t, router, http.MethodDelete, "/api/object/2", nil), http.StatusOK)
	// Removing a camera with history must fail without deleting its events.
	expectStatus(t, request(t, router, http.MethodDelete, "/api/camera/1", nil), http.StatusConflict)
	expectStatus(t, request(t, router, http.MethodGet, "/api/camera/1", nil), http.StatusOK)
	// The event keeps the referenced object from being deleted.
	expectStatus(t, request(t, router, http.MethodDelete, "/api/object/1", nil), http.StatusInternalServerError)
	expectStatus(t, request(t, router, http.MethodDelete, "/api/event/"+uuid, nil), http.StatusOK)
	expectStatus(t, request(t, router, http.MethodDelete, "/api/object/1", nil), http.StatusOK)
	for _, path := range []string{"/api/object/1", "/api/event/" + uuid} {
		for _, method := range []string{http.MethodGet, http.MethodDelete} {
			expectStatus(t, request(t, router, method, path, nil), http.StatusNotFound)
		}
	}
	expectStatus(t, request(t, router, http.MethodPut, "/api/object/1", entities.Object{Type: "missing"}), http.StatusNotFound)
	expectStatus(t, request(t, router, http.MethodPut, "/api/event/"+uuid, event), http.StatusNotFound)
	expectStatus(t, request(t, router, http.MethodDelete, "/api/camera/1", nil), http.StatusOK)
	expectStatus(t, request(t, router, http.MethodGet, "/api/camera/1", nil), http.StatusNotFound)
	expectStatus(t, request(t, router, http.MethodDelete, "/api/camera/1", nil), http.StatusNotFound)
	expectStatus(t, request(t, router, http.MethodPut, "/api/camera/1", map[string]any{"comment": "missing"}), http.StatusNotFound)
}

func TestEntityInvalidFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(nil).Router
	for _, entity := range []string{"camera", "object", "event"} {
		for _, action := range []string{"list", "count"} {
			for _, query := range []string{"camera_id=bad", "object_id=bad", "is_enabled=bad", "confidence_min=NaN", "confidence_min=1&confidence_max=0", "created_from=bad", "created_from=2026-11-02&created_to=2026-11-01"} {
				path := "/api/" + entity + "/" + action + "?" + query
				t.Run(path, func(t *testing.T) {
					expectStatus(t, request(t, router, http.MethodGet, path, nil), http.StatusBadRequest)
				})
			}
		}
	}
}

func TestCameraInvalidUpdates(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(nil).Router
	for _, body := range []any{json.RawMessage(`null`), map[string]any{}, map[string]any{"ip": "   "}, map[string]any{"ip": ""}, map[string]any{"port": 0}, map[string]any{"port": 65536}, map[string]any{"port": "bad"}, map[string]any{"is_enabled": "bad"}} {
		expectStatus(t, request(t, router, http.MethodPut, "/api/camera/1", body), http.StatusBadRequest)
	}
}
