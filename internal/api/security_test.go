package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityHeaders(t *testing.T) {
	h := NewHandler()
	for _, tc := range []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/"},
		{http.MethodGet, "/healthz"},
		{http.MethodPost, "/api/v1/logos/inspect"},
		{http.MethodPost, "/api/v1/favicons"},
	} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
		for key, want := range map[string]string{
			"X-Content-Type-Options": "nosniff",
			"Referrer-Policy": "no-referrer",
			"X-Frame-Options": "DENY",
			"Permissions-Policy": "camera=(), microphone=(), geolocation=()",
		} {
			if got := rec.Header().Get(key); got != want {
				t.Errorf("%s %s: %s = %q, want %q", tc.method, tc.path, key, got, want)
			}
		}
		if tc.method == http.MethodPost && rec.Header().Get("Cache-Control") != "no-store" {
			t.Errorf("%s: API response must not be cached", tc.path)
		}
	}
}
