package main

import (
	"InfoCenter/BackendServer"
	"log"
	"net/http"
)

// Middleware for CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", BackendServer.HomeHandler)
	mux.HandleFunc("/infocenter/", BackendServer.HandleRequests)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(mux)))
}
