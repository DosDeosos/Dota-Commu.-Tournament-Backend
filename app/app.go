package app

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func New() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", handleHealth)
	mux.HandleFunc("GET /api/data", handleData)
	return cors(mux)
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	if _, err := loadData(); err != nil {
		log.Printf("health: %v", err)
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"ok": false, "data": "invalid"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "data": "loaded"})
}

func handleData(w http.ResponseWriter, _ *http.Request) {
	snap, err := loadData()
	if err != nil {
		log.Printf("load data: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	w.Header().Set("Cache-Control", "public, s-maxage=60, stale-while-revalidate=300")
	writeJSON(w, http.StatusOK, snap)
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := os.Getenv("CORS_ORIGIN")
		if origin == "" {
			origin = "*"
		}
		h := w.Header()
		h.Set("Access-Control-Allow-Origin", origin)
		h.Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
