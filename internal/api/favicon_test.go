package api

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFaviconZIP(t *testing.T) {
	var input bytes.Buffer
	src := image.NewRGBA(image.Rect(0, 0, 40, 20))
	if err := png.Encode(&input, src); err != nil {
		t.Fatal(err)
	}
	req := uploadRequest(t, input.Bytes())
	req.URL.Path = "/api/v1/favicons"
	rec := httptest.NewRecorder()
	NewHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}
	z, err := zip.NewReader(bytes.NewReader(rec.Body.Bytes()), int64(rec.Body.Len()))
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]int{
		"favicon-16x16.png": 16,
		"favicon-32x32.png": 32,
		"apple-touch-icon.png": 180,
		"android-chrome-192x192.png": 192,
		"android-chrome-512x512.png": 512,
	}
	for _, file := range z.File {
		if size, ok := expected[file.Name]; ok {
			f, err := file.Open()
			if err != nil {
				t.Fatal(err)
			}
			cfg, err := png.DecodeConfig(f)
			_ = f.Close()
			if err != nil || cfg.Width != size || cfg.Height != size {
				t.Fatalf("%s dimensions: %+v, %v", file.Name, cfg, err)
			}
			delete(expected, file.Name)
		}
	}
	if len(expected) != 0 {
		t.Fatalf("missing assets: %v", expected)
	}
}

func TestICOHeader(t *testing.T) {
	ico := makeICO([][]byte{{1, 2}, {3, 4}}, []int{16, 32})
	if binary.LittleEndian.Uint16(ico[2:4]) != 1 || binary.LittleEndian.Uint16(ico[4:6]) != 2 {
		t.Fatal("invalid ICO directory")
	}
}

func TestRejectSVGGeneration(t *testing.T) {
	req := uploadRequest(t, []byte(`<svg viewBox="0 0 16 16"></svg>`))
	req.URL.Path = "/api/v1/favicons"
	rec := httptest.NewRecorder()
	NewHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("status %d", rec.Code)
	}
}
