package api

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"image/color"
	"image/draw"
	"strconv"
	"strings"
	"io"
	"net/http"

	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/webp"
)

const maxArchiveBytes = 32 << 20

func generateFavicons(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxLogoBytes+4096)
	if err := r.ParseMultipartForm(maxLogoBytes); err != nil {
		http.Error(w, "invalid multipart upload", http.StatusRequestEntityTooLarge)
		return
	}
	defer r.MultipartForm.RemoveAll()
	files := r.MultipartForm.File["logo"]
	if len(files) != 1 || files[0].Size > maxLogoBytes {
		http.Error(w, "one logo under 5 MiB required", http.StatusBadRequest)
		return
	}
	f, err := files[0].Open()
	if err != nil {
		http.Error(w, "invalid file", http.StatusBadRequest)
		return
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, maxLogoBytes+1))
	if err != nil || len(data) > maxLogoBytes {
		http.Error(w, "file too large", http.StatusRequestEntityTooLarge)
		return
	}
	var img image.Image
	if len(data) >= 8 && bytes.Equal(data[:8], []byte{137, 80, 78, 71, 13, 10, 26, 10}) {
		cfg, err := png.DecodeConfig(bytes.NewReader(data))
		if err != nil || !validDimensions(cfg.Width, cfg.Height) {
			http.Error(w, "invalid PNG dimensions", http.StatusBadRequest)
			return
		}
		img, err = png.Decode(bytes.NewReader(data))
	} else if len(data) >= 12 && bytes.Equal(data[:4], []byte("RIFF")) && bytes.Equal(data[8:12], []byte("WEBP")) {
		cfg, err := webp.DecodeConfig(bytes.NewReader(data))
		if err != nil || !validDimensions(cfg.Width, cfg.Height) {
			http.Error(w, "invalid WebP dimensions", http.StatusBadRequest)
			return
		}
		img, err = webp.Decode(bytes.NewReader(data))
	} else {
		http.Error(w, "PNG or WebP required; SVG rasterization is not available yet", http.StatusUnsupportedMediaType)
		return
	}
	if err != nil {
		http.Error(w, "invalid image", http.StatusBadRequest)
		return
	}
	options, err := parseIconOptions(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var archive bytes.Buffer
	zw := zip.NewWriter(&archive)
	sizes := []struct {
		name string
		size int
	}{
		{"favicon-16x16.png", 16},
		{"favicon-32x32.png", 32},
		{"apple-touch-icon.png", 180},
		{"android-chrome-192x192.png", 192},
		{"android-chrome-512x512.png", 512},
	}
	icoPNGs := make([][]byte, 0, 3)
	for _, item := range sizes {
		icon := renderSquareWithOptions(img, item.size, options)
		var out bytes.Buffer
		if err := png.Encode(&out, icon); err != nil {
			http.Error(w, "encoding failed", http.StatusInternalServerError)
			return
		}
		if err := addZipFile(zw, item.name, out.Bytes()); err != nil {
			http.Error(w, "archive failed", http.StatusInternalServerError)
			return
		}
		if item.size <= 32 {
			icoPNGs = append(icoPNGs, out.Bytes())
		}
	}
	ico48 := renderSquareWithOptions(img, 48, options)
	var out48 bytes.Buffer
	if err := png.Encode(&out48, ico48); err != nil {
		http.Error(w, "encoding failed", http.StatusInternalServerError)
		return
	}
	icoPNGs = append(icoPNGs, out48.Bytes())
	if err := addZipFile(zw, "favicon.ico", makeICO(icoPNGs, []int{16, 32, 48})); err != nil {
		http.Error(w, "archive failed", http.StatusInternalServerError)
		return
	}
	manifest, _ := json.MarshalIndent(map[string]any{
		"name":       "My Website",
		"short_name": "Website",
		"icons": []map[string]string{
			{"src": "/android-chrome-192x192.png", "sizes": "192x192", "type": "image/png"},
			{"src": "/android-chrome-512x512.png", "sizes": "512x512", "type": "image/png"},
		},
		"display": "standalone",
	}, "", "  ")
	if err := addZipFile(zw, "site.webmanifest", manifest); err != nil {
		http.Error(w, "archive failed", http.StatusInternalServerError)
		return
	}
	snippet := []byte(`<link rel="icon" href="/favicon.ico" sizes="any">
<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
<link rel="apple-touch-icon" href="/apple-touch-icon.png">
<link rel="manifest" href="/site.webmanifest">
`)
	if err := addZipFile(zw, "head.html", snippet); err != nil {
		http.Error(w, "archive failed", http.StatusInternalServerError)
		return
	}
	if err := zw.Close(); err != nil || archive.Len() > maxArchiveBytes {
		http.Error(w, "archive too large", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\"kaqen-favicons.zip\"")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = w.Write(archive.Bytes())
}

type iconOptions struct {
	Padding int
	Background color.NRGBA
	UseBackground bool
}

func parseIconOptions(r *http.Request) (iconOptions, error) {
	var opt iconOptions
	raw := r.FormValue("padding")
	if raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value < 0 || value > 40 {
			return opt, fmt.Errorf("padding must be between 0 and 40 percent")
		}
		opt.Padding = value
	}
	background := strings.TrimSpace(r.FormValue("background"))
	if background != "" && background != "transparent" {
		if len(background) != 7 || background[0] != '#' {
			return opt, fmt.Errorf("background must be #RRGGBB or transparent")
		}
		value, err := strconv.ParseUint(background[1:], 16, 24)
		if err != nil {
			return opt, fmt.Errorf("background must be #RRGGBB or transparent")
		}
		opt.Background = color.NRGBA{R: uint8(value >> 16), G: uint8(value >> 8), B: uint8(value), A: 255}
		opt.UseBackground = true
	}
	return opt, nil
}

func renderSquareWithOptions(src image.Image, size int, opt iconOptions) *image.NRGBA {
	dst := image.NewNRGBA(image.Rect(0, 0, size, size))
	if opt.UseBackground {
		draw.Draw(dst, dst.Bounds(), &image.Uniform{C: opt.Background}, image.Point{}, draw.Src)
	}
	bounds := src.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	if width < 1 || height < 1 { return dst }
	available := size * (100 - 2*opt.Padding) / 100
	if available < 1 { available = 1 }
	targetW, targetH := available, available
	if width > height {
		targetH = max(1, int(float64(available)*float64(height)/float64(width)))
	} else {
		targetW = max(1, int(float64(available)*float64(width)/float64(height)))
	}
	rect := image.Rect((size-targetW)/2, (size-targetH)/2, (size-targetW)/2+targetW, (size-targetH)/2+targetH)
	xdraw.CatmullRom.Scale(dst, rect, src, bounds, xdraw.Over, nil)
	return dst
}

func validDimensions(width, height int) bool {
	return width > 0 && height > 0 && width <= maxPixels/height
}

func renderSquare(src image.Image, size int) *image.NRGBA {
	dst := image.NewNRGBA(image.Rect(0, 0, size, size))
	bounds := src.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	if width < 1 || height < 1 {
		return dst
	}
	targetW, targetH := size, size
	if width > height {
		targetH = max(1, int(float64(size)*float64(height)/float64(width)))
	} else {
		targetW = max(1, int(float64(size)*float64(width)/float64(height)))
	}
	rect := image.Rect((size-targetW)/2, (size-targetH)/2, (size-targetW)/2+targetW, (size-targetH)/2+targetH)
	xdraw.CatmullRom.Scale(dst, rect, src, bounds, xdraw.Src, nil)
	return dst
}

func addZipFile(z *zip.Writer, name string, data []byte) error {
	f, err := z.Create(name)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}

// ICO supports PNG-compressed images in directory entries on modern browsers.
func makeICO(images [][]byte, sizes []int) []byte {
	var b bytes.Buffer
	_ = binary.Write(&b, binary.LittleEndian, uint16(0))
	_ = binary.Write(&b, binary.LittleEndian, uint16(1))
	_ = binary.Write(&b, binary.LittleEndian, uint16(len(images)))
	offset := 6 + 16*len(images)
	for i, data := range images {
		b.WriteByte(byte(sizes[i]))
		b.WriteByte(byte(sizes[i]))
		b.WriteByte(0)
		b.WriteByte(0)
		_ = binary.Write(&b, binary.LittleEndian, uint16(1))
		_ = binary.Write(&b, binary.LittleEndian, uint16(32))
		_ = binary.Write(&b, binary.LittleEndian, uint32(len(data)))
		_ = binary.Write(&b, binary.LittleEndian, uint32(offset))
		offset += len(data)
	}
	for _, data := range images {
		b.Write(data)
	}
	return b.Bytes()
}
