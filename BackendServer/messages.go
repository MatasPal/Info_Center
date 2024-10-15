package BackendServer

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var messageID int
var mu sync.Mutex

// Handles posting a message to a topic
func handlePost(w http.ResponseWriter, r *http.Request, topicName string) {
	// Allow CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Your existing logic for handling POST requests
	var msgData string
	fmt.Fscan(r.Body, &msgData)
	r.Body.Close()

	mu.Lock()
	messageID++
	msg := Message{ID: messageID, Data: msgData}
	mu.Unlock()

	log.Printf("Message received: %s", msg.Data)

	topic := getTopic(topicName)
	topic.mu.Lock()
	topic.Messages = append(topic.Messages, msg)

	for client, ch := range topic.Clients {
		select {
		case ch <- msg:
			log.Printf("Message sent to client: %s", msg.Data)
		default:
			close(ch)
			delete(topic.Clients, client)
			log.Printf("Client disconnected")
		}
	}
	topic.mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func handleGet(w http.ResponseWriter, r *http.Request, topicName string) {
	// Allow CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Your existing logic for handling GET requests (Server-Sent Events)
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	topic := getTopic(topicName)
	ch := make(chan Message)
	topic.mu.Lock()
	topic.Clients[w] = ch

	// Send all previously sent messages to the new client
	for _, msg := range topic.Messages {
		fmt.Fprintf(w, "id: %d\nevent: msg\ndata: %s\n\n", msg.ID, msg.Data)
		flusher.Flush()
	}
	topic.mu.Unlock()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	timeout := time.After(30 * time.Second)
	start := time.Now()

	log.Printf("Client subscribed to topic: %s", topicName)

	for {
		select {
		case msg := <-ch:
			log.Printf("Sending message to client: %s", msg.Data)
			fmt.Fprintf(w, "id: %d\nevent: msg\ndata: %s\n\n", msg.ID, msg.Data)
			flusher.Flush()
		case <-timeout:
			duration := time.Since(start)
			log.Printf("Client disconnected after %d seconds", int(duration.Seconds()))
			fmt.Fprintf(w, "event: timeout\ndata: %ds\n\n", int(duration.Seconds()))
			flusher.Flush()
			topic.mu.Lock()
			delete(topic.Clients, w)
			topic.mu.Unlock()
			removeTopicIfEmpty(topicName)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
