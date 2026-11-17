package api

import (
	"archive/zip"
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEnhancedExportAssets(t *testing.T) {
	var input bytes.Buffer
	src := image.NewNRGBA(image.Rect(0, 0, 64, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			src.SetNRGBA(x, y, color.NRGBA{R: 255, A: 255})
		}
	}
	if err := png.Encode(&input, src); err != nil {
		t.Fatal(err)
	}
	req := uploadRequest(t, input.Bytes())
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
	files := make(map[string][]byte)
	for _, file := range z.File {
		if file.Name != "og-image.png" && file.Name != "README.txt" && file.Name != "head.html" {
			continue
		}
		f, err := file.Open()
		if err != nil {
			t.Fatal(err)
		}
		data, err := io.ReadAll(f)
		_ = f.Close()
		if err != nil {
			t.Fatal(err)
		}
		files[file.Name] = data
	}
	cfg, err := png.DecodeConfig(bytes.NewReader(files["og-image.png"]))
	if err != nil || cfg.Width != 1200 || cfg.Height != 630 {
		t.Fatalf("invalid social card dimensions: %+v, %v", cfg, err)
	}
	for _, token := range []string{"og:image", "twitter:card", "og-image.png"} {
		if !strings.Contains(string(files["head.html"]), token) {
			t.Fatalf("missing integration token %q", token)
		}
	}
	if !strings.Contains(string(files["README.txt"]), "https://example.com") {
		t.Fatal("missing setup instructions")
	}
}

func TestSocialCardAspectRatioAndOpaqueBackground(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 100, 50))
	for y := 0; y < 50; y++ {
		for x := 0; x < 100; x++ {
			src.SetNRGBA(x, y, color.NRGBA{R: 255, A: 255})
		}
	}
	card := renderSocialCard(src, iconOptions{})
	if corner := card.NRGBAAt(0, 0); corner != (color.NRGBA{R: 255, G: 255, B: 255, A: 255}) {
		t.Fatalf("unexpected background: %+v", corner)
	}
	if center := card.NRGBAAt(600, 315); center.R != 255 || center.G != 0 || center.A != 255 {
		t.Fatalf("logo missing from center: %+v", center)
	}
}
