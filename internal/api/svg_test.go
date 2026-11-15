package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSVGInspect(t *testing.T) {
	valid := []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 32"><rect width="64" height="32" fill="#000"/></svg>`)
	w := httptest.NewRecorder()
	NewHandler().ServeHTTP(w, uploadRequest(t, valid))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"format":"svg"`) {
		t.Fatalf("unexpected response: %d %s", w.Code, w.Body.String())
	}
}

func TestSVGRejectUnsafe(t *testing.T) {
	inputs := []string{
		`<svg viewBox="0 0 10 10"><script>alert(1)</script></svg>`,
		`<svg viewBox="0 0 10 10"><image href="https://example.org/x"/></svg>`,
		`<!DOCTYPE svg [<!ENTITY x SYSTEM "file:///etc/passwd">]><svg viewBox="0 0 10 10"/>`,
		`<svg viewBox="0 0 10 10"><rect style="fill:url(https://example.org/x)"/></svg>`,
	}
	for _, input := range inputs {
		w := httptest.NewRecorder()
		NewHandler().ServeHTTP(w, uploadRequest(t, []byte(input)))
		if w.Code == http.StatusOK {
			t.Fatalf("unsafe SVG accepted: %s", input)
		}
	}
}
