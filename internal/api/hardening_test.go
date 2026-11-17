package api

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestHealthAndReadiness(t *testing.T) {
	h := NewHandler()
	for _, path := range []string{"/healthz", "/readyz"} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "status") {
			t.Fatalf("%s: status %d, body %s", path, rec.Code, rec.Body.String())
		}
	}
}

func TestConversionCapacity(t *testing.T) {
	h := NewHandler()
	for i := 0; i < cap(conversionSlots); i++ {
		conversionSlots <- struct{}{}
	}
	defer func() {
		for i := 0; i < cap(conversionSlots); i++ {
			<-conversionSlots
		}
	}()
	for _, path := range []string{"/api/v1/logos/inspect", "/api/v1/favicons"} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, path, nil))
		if rec.Code != http.StatusServiceUnavailable {
			t.Fatalf("%s: expected 503, got %d", path, rec.Code)
		}
		if rec.Header().Get("Retry-After") != strconv.Itoa(5) {
			t.Fatalf("%s: missing retry hint", path)
		}
	}
}
