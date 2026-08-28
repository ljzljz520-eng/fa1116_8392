package api

import (
	"gestureparticles/internal/flow021"
	"gestureparticles/internal/importer"
	"gestureparticles/internal/registry"
	"gestureparticles/internal/review"
	"gestureparticles/internal/store"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthAndRegisterEndpoint(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/api.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	s := registry.New(db, review.New(db), importer.New(db), flow021.New())
	h := New(s)
	request := httptest.NewRequest("GET", "/health", nil)
	response := httptest.NewRecorder()
	h.ServeHTTP(response, request)
	if response.Code != 200 || !strings.Contains(response.Body.String(), "ok") {
		t.Fatalf("code=%d body=%s", response.Code, response.Body.String())
	}
}
