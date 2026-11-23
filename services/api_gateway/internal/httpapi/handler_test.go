package httpapi

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLivenessDoesNotDependOnDatabase(t *testing.T) {
	h := Handler(func(context.Context) error { t.Fatal("liveness queried database"); return nil }, slog.New(slog.NewTextHandler(io.Discard, nil)), nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if w.Code != 200 {
		t.Fatal(w.Code)
	}
}

func TestReadinessAndErrors(t *testing.T) {
	for _, ready := range []bool{true, false} {
		h := Handler(func(ctx context.Context) error {
			if _, ok := ctx.Deadline(); !ok {
				t.Error("no query timeout")
			}
			if !ready {
				return errors.New("internal database credentials should never appear in the response")
			}
			return nil
		}, slog.New(slog.NewTextHandler(io.Discard, nil)), nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/readyz", nil))
		want := 503
		if ready {
			want = 200
		}
		if w.Code != want {
			t.Fatal(w.Code)
		}
		if strings.Contains(w.Body.String(), "credentials") {
			t.Fatal("internal error exposed")
		}
		if w.Header().Get("Cache-Control") != "no-store" {
			t.Fatal("health response can be cached")
		}
	}
}

func TestNoCrudRoutesYet(t *testing.T) {
	h := Handler(func(context.Context) error { return nil }, slog.New(slog.NewTextHandler(io.Discard, nil)), nil)
	for _, path := range []string{"/", "/api/objects", "/unknown"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		if w.Code != 404 {
			t.Fatalf("%s: %d", path, w.Code)
		}
	}
}
