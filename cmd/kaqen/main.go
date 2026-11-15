package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Ploos-AS/Kaqen/internal/api"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	srv := &http.Server{Addr: addr, Handler: api.NewHandler(), ReadHeaderTimeout: 5 * time.Second}
	log.Printf("Kaqen listening on %s", addr)
	log.Fatal(srv.ListenAndServe())
}
