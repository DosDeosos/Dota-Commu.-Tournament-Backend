package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DosDeosos/dota-commu-tournament-backend/app"
)

func main() {
	addr := ":" + envOr("PORT", "8080")
	log.Printf("listening on http://localhost%s", addr)
	srv := &http.Server{
		Addr:              addr,
		Handler:           app.New(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
