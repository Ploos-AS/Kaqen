package api

import (
	"archive/zip"
	"bytes"
	"image/png"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRasterizeStaticSVG(t *testing.T) {
	data := []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64"><rect x="0" y="0" width="64" height="64" fill="#ff0000"/></svg>`)
	img, err := rasterizeSVG(data)
	if err != nil {
		t.Fatal(err)
	}
	if img.Bounds().Dx() != 64 || img.Bounds().Dy() != 64 {
		t.Fatalf("wrong SVG dimensions: %v", img.Bounds())
	}
}

func TestSVGRejected(t *testing.T) {
	for _, svg := range []string{
		`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64"><script>alert(1)</script></svg>`,
		`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64"><image href="https://example.com/a.png"/></svg>`,
		`<!DOCTYPE svg [<!ENTITY x SYSTEM "file:///etc/passwd">]><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64"/>`,
	} {
		if _, err := rasterizeSVG([]byte(svg)); err == nil {
			t.Fatalf("unsafe SVG accepted: %s", svg)
		}
	}
}

func TestSVGFaviconZIP(t *testing.T) {
	data := []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64"><circle cx="32" cy="32" r="20" fill="#00ff00"/></svg>`)
	req := uploadRequest(t, data)
	req.URL.Path = "/api/v1/favicons"
	rec := httptest.NewRecorder()
	NewHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("HTTP %d: %s", rec.Code, rec.Body.String())
	}
	archive, err := zip.NewReader(bytes.NewReader(rec.Body.Bytes()), int64(rec.Body.Len()))
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range archive.File {
		if f.Name == "favicon-32x32.png" {
			rc, err := f.Open()
			if err != nil {
				t.Fatal(err)
			}
			defer rc.Close()
			cfg, err := png.DecodeConfig(rc)
			if err != nil || cfg.Width != 32 {
				t.Fatalf("bad favicon: %v, %v", cfg, err)
			}
			return
		}
	}
	t.Fatal("favicon missing")
}

func TestSVGComplexityLimits(t *testing.T) {
	prefix := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">`
	for _, body := range []string{
		strings.Repeat("<rect width=\"1\" height=\"1\"/>", 2050),
		"<path d=\"" + strings.Repeat("M1 1 ", 7000) + "\"/>",
	} {
		if _, err := rasterizeSVG([]byte(prefix + body + "</svg>")); err == nil {
			t.Fatal("excessively complex SVG accepted")
		}
	}
}
