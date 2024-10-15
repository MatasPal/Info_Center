package BackendServer

import (
	"fmt"
	"net/http"
)

// Exported function: Handles root "/"
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Welcome to the InfoCenter!")
}

// Exported function: Handles "/infocenter/{topic}" requests
func HandleRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	topicName := r.URL.Path[len("/infocenter/"):]
	if topicName == "" {
		http.Error(w, "Topic not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodPost:
		handlePost(w, r, topicName)
	case http.MethodGet:
		handleGet(w, r, topicName)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
