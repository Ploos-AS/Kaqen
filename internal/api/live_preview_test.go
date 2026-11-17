package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLivePreviewMarkup(t *testing.T) {
	rec := httptest.NewRecorder()
	NewHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("homepage status %d", rec.Code)
	}
	html := rec.Body.String()
	for _, expected := range []string{
		"id=\"output-preview\"",
		"id=\"custom-color\"",
		"data-icon-size=\"16\"",
		"data-icon-size=\"32\"",
		"data-icon-size=\"180\"",
		"data-icon-size=\"512\"",
	} {
		if !strings.Contains(html, expected) {
			t.Errorf("missing %s", expected)
		}
	}
}
