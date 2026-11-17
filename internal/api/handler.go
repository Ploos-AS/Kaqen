package api

import (
	"embed"
	"io/fs"
	"encoding/json"
	"net/http"
)

//go:embed web/*
var webAssets embed.FS

//go:embed web/index.html
var indexHTML []byte

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	assets, err := fs.Sub(webAssets, "web")
	if err != nil {
		panic(err)
	}
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assets))))
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self' blob:; connect-src 'self'; object-src 'none'; base-uri 'none'; form-action 'self'")
		_, _ = w.Write(indexHTML)
	})
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	mux.HandleFunc("POST /api/v1/logos/inspect", uploadLogo)
	mux.HandleFunc("POST /api/v1/favicons", generateFavicons)
	return mux
}
