package api

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"io"	"net/http/httptest"
	"testing"
)

func TestMaskableSafeZone(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			src.SetNRGBA(x, y, color.NRGBA{R: 255, A: 255})
		}
	}
	bg := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	out := renderMaskable(src, 192, bg)
	if out.Bounds().Dx() != 192 || out.Bounds().Dy() != 192 {
		t.Fatal("wrong dimensions")
	}
	if got := out.NRGBAAt(0, 0); got != bg {
		t.Fatalf("corner must be opaque background: %+v", got)
	}
	if got := out.NRGBAAt(96, 96); got.R != 255 || got.G != 0 {
		t.Fatalf("center must contain logo: %+v", got)
	}
	// All logo pixels must lie inside the central safe-zone circle.
	for y := 0; y < 192; y++ {
		for x := 0; x < 192; x++ {
			p := out.NRGBAAt(x, y)
			if p.G == 0 && p.B == 0 && (x-96)*(x-96)+(y-96)*(y-96) > 76*76 {
				t.Fatalf("logo outside safe zone at %d,%d", x, y)
			}
		}
	}
}

func TestMaskableArchiveAndManifest(t *testing.T) {
	var input bytes.Buffer
	src := image.NewNRGBA(image.Rect(0, 0, 64, 64))
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
	expected := map[string]int{
		"android-chrome-maskable-192x192.png": 192,
		"android-chrome-maskable-512x512.png": 512,
	}
	manifestFound := false
	for _, f := range z.File {
		reader, err := f.Open()
		if err != nil {
			t.Fatal(err)
		}
		if size, ok := expected[f.Name]; ok {
			cfg, err := png.DecodeConfig(reader)
			_ = reader.Close()
			if err != nil || cfg.Width != size || cfg.Height != size {
				t.Fatalf("invalid %s: %+v %v", f.Name, cfg, err)
			}
			delete(expected, f.Name)
			continue
		}
		if f.Name == "site.webmanifest" {
			var manifest struct {
				Icons []struct {
					Src string `json:"src"`
					Purpose string `json:"purpose"`
				} `json:"icons"`
			}
			data, err := io.ReadAll(reader)
			if err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(data, &manifest); err != nil {
				t.Fatal(err)
			}
			count := 0
			for _, icon := range manifest.Icons {
				if icon.Purpose == "maskable" {
					count++
				}
			}
			if count != 2 {
				t.Fatalf("expected 2 maskable icons, got %d", count)
			}
			manifestFound = true
		}
		_ = reader.Close()
	}
	if len(expected) != 0 || !manifestFound {
		t.Fatalf("missing assets %v, manifest=%v", expected, manifestFound)
	}
}
