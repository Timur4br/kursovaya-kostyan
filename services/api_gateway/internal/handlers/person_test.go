package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

type personControllerStub struct {
	calls   int
	failure error
	input   dto.PersonInput
}

func (s *personControllerStub) Create(_ context.Context, input dto.PersonInput) (models.Person, error) {
	s.calls++
	s.input = input
	return models.Person{ID: 1}, s.failure
}
func (s *personControllerStub) Get(context.Context, int64) (models.Person, error) {
	s.calls++
	return models.Person{ID: 1}, s.failure
}
func (s *personControllerStub) Update(_ context.Context, _ int64, input dto.PersonInput) (models.Person, error) {
	s.calls++
	s.input = input
	return models.Person{ID: 1}, s.failure
}
func (s *personControllerStub) Delete(context.Context, int64) error { s.calls++; return s.failure }
func personRouter(s *personControllerStub) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewPersonHandler(s, slog.New(slog.NewTextHandler(io.Discard, nil)))
	r.POST("/person", h.Create)
	r.GET("/person/:id", h.Get)
	r.PUT("/person/:id", h.Update)
	r.DELETE("/person/:id", h.Delete)
	return r
}

func TestPersonHandlerRejectsInvalidTransportInput(t *testing.T) {
	cases := []struct {
		method, path, body, media string
		status                    int
	}{
		{"GET", "/person/0", "", "", 400},
		{"GET", "/person/9223372036854775808", "", "", 400},
		{"POST", "/person", `[]`, "application/json", 400},
		{"POST", "/person", `{} {}`, "application/json", 400},
		{"POST", "/person", `{"user_id":1,"user_id":2}`, "application/json", 400},
		{"POST", "/person", `{`, "application/json", 400},
		{"POST", "/person", `{}`, "text/plain", 415},
		{"PUT", "/person/1", `{"id":2}`, "application/json", 422},
		{"POST", "/person", `{"user_id":"1"}`, "application/json", 422},
		{"POST", "/person", `{"unknown":true}`, "application/json", 422},
		{"POST", "/person", `{"visual_features":{"large":"` + strings.Repeat("a", int(MaxBodyBytes)) + `"}}`, "application/json", 413},
	}
	for _, test := range cases {
		s := &personControllerStub{}
		r := personRouter(s)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
		if test.media != "" {
			req.Header.Set("Content-Type", test.media)
		}
		r.ServeHTTP(w, req)
		if w.Code != test.status {
			t.Fatalf("%s %s: %d %s", test.method, test.path, w.Code, w.Body)
		}
		if s.calls != 0 {
			t.Fatal("invalid transport reached controller")
		}
		if !json.Valid(w.Body.Bytes()) {
			t.Fatal("error is not JSON")
		}
	}
}

func TestPersonHandlerUsesControllerAndReturnsStatus(t *testing.T) {
	s := &personControllerStub{}
	r := personRouter(s)
	for _, method := range []string{"POST", "PUT"} {
		path := "/person"
		if method == "PUT" {
			path += "/1"
		}
		req := httptest.NewRequest(method, path, strings.NewReader(`{"user_id":null}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		want := 201
		if method == "PUT" {
			want = 200
		}
		if w.Code != want || w.Header().Get("Location") != "/person/1" {
			t.Fatal(w.Code, w.Header(), w.Body)
		}
		if !s.input.UserID.Present || !s.input.UserID.Null {
			t.Fatal("explicit null lost")
		}
	}
	for _, test := range []struct {
		err    error
		status int
	}{
		{sql.ErrNoRows, 404},
		{&pgconn.PgError{Code: "23503", Detail: "private SQL detail"}, 409},
		{&pgconn.PgError{Code: "23514", Detail: "private SQL detail"}, 422},
		{context.DeadlineExceeded, 504},
	} {
		s.failure = test.err
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("GET", "/person/1", nil))
		if w.Code != test.status || strings.Contains(w.Body.String(), "private") {
			t.Fatal(w.Code, w.Body)
		}
	}
}
