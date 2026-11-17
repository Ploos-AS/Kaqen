package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWebHome(t *testing.T) {
	rec := httptest.NewRecorder()
	NewHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "Kaqen") {
		t.Fatalf("unexpected homepage: %d", rec.Code)
	}
	if rec.Header().Get("Content-Security-Policy") == "" {
		t.Fatal("missing Content Security Policy")
	}
}

func TestWebAssets(t *testing.T) {
	for _, path := range []string{"/assets/style.css", "/assets/app.js"} {
		rec := httptest.NewRecorder()
		NewHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: HTTP %d", path, rec.Code)
		}
	}
}

func TestWebUnknownRoute(t *testing.T) {
	rec := httptest.NewRecorder()
	NewHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/does-not-exist", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("HTTP %d", rec.Code)
	}
}
