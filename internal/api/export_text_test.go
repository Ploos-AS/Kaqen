package api

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExportTextUsesRealNewlines(t *testing.T) {
	pngData := samplePNG(t, 64, 64)
	req := uploadRequest(t, pngData)
	req.URL.Path = "/api/v1/favicons"
	rec := httptest.NewRecorder()
	NewHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("HTTP %d: %s", rec.Code, rec.Body.String())
	}
	z, err := zip.NewReader(bytes.NewReader(rec.Body.Bytes()), int64(rec.Body.Len()))
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"head.html", "README.txt"} {
		var found bool
		for _, f := range z.File {
			if f.Name != name {
				continue
			}
			found = true
			rc, err := f.Open()
			if err != nil {
				t.Fatal(err)
			}
			data, err := io.ReadAll(rc)
			_ = rc.Close()
			if err != nil {
				t.Fatal(err)
			}
			if bytes.Contains(data, []byte(`\n`)) {
				t.Fatalf("%s contains literal escaped newlines", name)
			}
			if !strings.Contains(string(data), "\n") {
				t.Fatalf("%s has no real line breaks", name)
			}
		}
		if !found {
			t.Fatalf("%s missing from ZIP", name)
		}
	}
}
