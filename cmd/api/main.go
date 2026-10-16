package main

import (
	"blog/internal/database"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := database.NewPool(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	mux.HandleFunc("GET /health/db", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{
				"status": "unavailable",
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("API Listening on %s", addr)
	log.Fatal(server.ListenAndServe())
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}
