package BackendServer

import (
	"net/http"
)

// HomeHandler - Handles root "/"
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Welcome to the InfoCenter!"))
}

// HandleRequests - Handles "/infocenter/{topic}" requests
func HandleRequests(w http.ResponseWriter, r *http.Request) {
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
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
