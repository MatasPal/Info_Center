package main

import (
	"InfoCenter/BackendServer"
	"log"
	"net/http"
)

// Initializes the HTTP server and defines the root handler.
func main() {
	// Set up HTTP handlers
	http.HandleFunc("/", BackendServer.HomeHandler)
	http.HandleFunc("/infocenter/", BackendServer.HandleRequests)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
