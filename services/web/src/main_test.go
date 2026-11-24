package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProxyForwardsGatewayRequests(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/event/list":
			if r.Method != http.MethodGet || r.URL.Query().Get("camera_id") != "7" || r.URL.Query().Get("offset") != "15" || r.URL.Query().Get("created_from") != "2026-11-01T00:00:00+03:00" {
				t.Error("filter or pagination was lost by proxy")
			}
			_, _ = w.Write([]byte(`[{"camera_id":7,"type":"found"}]`))
		case "/api/camera":
			if r.Method != http.MethodPost || r.Header.Get("Content-Type") != "application/json" {
				t.Error("incorrect create request")
			}
			var camera struct {
				IP        string
				Port      int
				IsEnabled bool `json:"is_enabled"`
			}
			if err := json.NewDecoder(r.Body).Decode(&camera); err != nil {
				t.Error(err)
			}
			if camera.IP != "192.168.1.10" || camera.Port != 80 || camera.IsEnabled {
				t.Error("camera body changed")
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`42`))
		case "/api/object/1":
			w.WriteHeader(http.StatusNotFound)
		default:
			t.Errorf("unexpected upstream request: %s", r.URL.Path)
		}
	}))
	defer upstream.Close()
	handler, err := newHandler(upstream.URL)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		method, path, body string
		status             int
		response           string
	}{
		{http.MethodGet, "/api/event/list?camera_id=7&offset=15&created_from=2026-11-01T00%3A00%3A00%2B03%3A00", "", http.StatusOK, `[{"camera_id":7,"type":"found"}]`},
		{http.MethodPost, "/api/camera", `{"ip":"192.168.1.10","port":80,"is_enabled":false}`, http.StatusCreated, "42"},
		{http.MethodGet, "/api/object/1", "", http.StatusNotFound, ""},
	} {
		t.Run(tc.method+tc.path, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, req)
			if response.Code != tc.status || response.Body.String() != tc.response {
				t.Fatalf("response = %d %s", response.Code, response.Body.String())
			}
		})
	}
}

func TestWebAssetsAndHealth(t *testing.T) {
	handler, err := newHandler("http://127.0.0.1:1")
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/", "/app.js", "/style.css", "/favicon.svg", "/health"} {
		t.Run(path, func(t *testing.T) {
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
			if response.Code != http.StatusOK {
				t.Fatalf("status = %d", response.Code)
			}
			if response.Header().Get("Content-Security-Policy") == "" {
				t.Fatal("missing content policy")
			}
		})
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/missing.js", nil))
	if response.Code != http.StatusNotFound {
		t.Fatalf("missing file: %d", response.Code)
	}
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/camera/list", nil))
	if response.Code != http.StatusBadGateway {
		t.Fatalf("unavailable gateway: %d", response.Code)
	}
	body, _ := io.ReadAll(response.Result().Body)
	if !json.Valid(body) {
		t.Fatal("gateway failure must return JSON")
	}
}

func TestInvalidGatewayURL(t *testing.T) {
	for _, target := range []string{"", "api_gateway:8080", "file:///tmp/db", "http://user:password@host", "http://host/api", "http://host?x=y", "http://host#fragment"} {
		if _, err := newHandler(target); err == nil {
			t.Errorf("accepted invalid target %q", target)
		}
	}
}
