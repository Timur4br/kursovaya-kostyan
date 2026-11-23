package httpapi

import (
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
	"github.com/gin-gonic/gin"
)

func TestExplicitEntityRoutesMatchFourOperationContract(t *testing.T) {
	r := Handler(nil, slog.New(slog.NewTextHandler(io.Discard, nil)), &Controllers{}).(*gin.Engine)
	registered := map[string]bool{}
	for _, route := range r.Routes() {
		registered[route.Method+" "+route.Path] = true
	}
	for _, resource := range resources.All {
		root := "/" + resource.Path
		item := root
		for _, key := range resource.Keys {
			item += "/:" + key
		}
		for _, route := range []string{"POST " + root, "GET " + item, "PUT " + item, "DELETE " + item} {
			if !registered[route] {
				t.Fatal("missing route", route)
			}
		}
		if registered["PATCH "+item] || registered["GET "+root] {
			t.Fatal("unexpected operation", resource.Path)
		}
	}
	for key := range registered {
		if strings.Contains(key, "/api/v1/") {
			t.Fatal("old generic API route retained")
		}
	}
	if len(registered) != 55 {
		t.Fatalf("want 52 entity and 3 utility routes, got %d", len(registered))
	}
	paths := OpenAPI()["paths"].(schemaMap)
	if len(paths) != 26 {
		t.Fatal("contract missing entity paths")
	}
}
