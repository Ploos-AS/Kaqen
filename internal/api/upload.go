package api

import (
	"bytes"
	"encoding/json"
	"image/png"
	"strings"
	"io"
	"net/http"
)

const maxLogoBytes = 5 << 20
const maxPixels = 16_000_000

type uploadResponse struct {
	Format string `json:"format"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Bytes  int    `json:"bytes"`
}

func uploadLogo(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxLogoBytes+4096)
	if err := r.ParseMultipartForm(maxLogoBytes); err != nil {
		http.Error(w, "invalid or oversized upload", http.StatusRequestEntityTooLarge)
		return
	}
	defer r.MultipartForm.RemoveAll()

	files := r.MultipartForm.File["logo"]
	if len(files) != 1 {
		http.Error(w, "exactly one logo required", http.StatusBadRequest)
		return
	}
	if files[0].Size > maxLogoBytes {
		http.Error(w, "logo too large", http.StatusRequestEntityTooLarge)
		return
	}
	f, err := files[0].Open()
	if err != nil {
		http.Error(w, "cannot open file", http.StatusBadRequest)
		return
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, maxLogoBytes+1))
	if err != nil || len(data) > maxLogoBytes {
		http.Error(w, "logo too large", http.StatusRequestEntityTooLarge)
		return
	}
	if strings.HasPrefix(strings.TrimSpace(string(data)), "<svg") || strings.HasPrefix(strings.TrimSpace(string(data)), "<?xml") {
		width, height, err := inspectSVG(data)
		if err != nil {
			http.Error(w, "unsupported or unsafe SVG", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(uploadResponse{"svg", width, height, len(data)})
		return
	}
	if len(data) < 8 || !bytes.Equal(data[:8], []byte{137, 80, 78, 71, 13, 10, 26, 10}) {
		http.Error(w, "PNG or restricted SVG only in this milestone", http.StatusUnsupportedMediaType)
		return
	}
	cfg, err := png.DecodeConfig(bytes.NewReader(data))
	if err != nil || cfg.Width < 1 || cfg.Height < 1 || cfg.Width > maxPixels/cfg.Height {
		http.Error(w, "invalid dimensions", http.StatusBadRequest)
		return
	}
	if _, err = png.Decode(bytes.NewReader(data)); err != nil {
		http.Error(w, "corrupt PNG", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(uploadResponse{"png", cfg.Width, cfg.Height, len(data)})
}
