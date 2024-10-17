package BackendServer

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Handles posting a message to a topic
func handlePost(w http.ResponseWriter, r *http.Request, topicName string) {
	// Read message data from request body
	msgData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read message data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Generate and append the new message
	msgID := getNextMessageID()
	msg := Message{ID: msgID, Data: string(msgData)}

	topic := getTopic(topicName)
	topic.mu.Lock()
	topic.Messages = append(topic.Messages, msg)

	// Broadcast the message to all connected clients
	for client, ch := range topic.Clients {
		select {
		case ch <- msg: // Send the message
		default:
			close(ch) // Disconnect slow clients
			delete(topic.Clients, client)
			log.Printf("Disconnected slow client from topic %s", topicName)
			removeTopicIfEmpty(topicName) // Check if topic should be removed
		}
	}
	topic.mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

// Handles retrieving messages for a topic
func handleGet(w http.ResponseWriter, r *http.Request, topicName string) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Check if topic exists
	topic := getTopic(topicName)
	if topic == nil {
		http.Error(w, "Topic not found", http.StatusNotFound)
		return
	}

	clientChan := make(chan Message, 10) // Buffered channel for slow clients
	topic.mu.Lock()
	topic.Clients[w] = clientChan
	topic.mu.Unlock()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Track connection time for timeout
	startTime := time.Now()
	timeoutDuration := 30 * time.Second // Set the timeout duration

	for {
		select {
		case msg := <-clientChan:
			if msg.Data != "" { // Avoid empty messages
				// Send message event in the required format
				fmt.Fprintf(w, "id: %d\nevent: msg\ndata: %s\n\n", msg.ID, msg.Data)
				flusher.Flush()
			}
		case <-time.After(timeoutDuration): // Timeout event after 30 seconds
			// Calculate the total time the client was connected (30s + 1s)
			duration := int(time.Since(startTime).Seconds()) + 1

			// Send the timeout event in the required format
			fmt.Fprintf(w, "event: timeout\ndata: %ds\n\n", duration)
			flusher.Flush()

			// Clean up client and check if topic can be removed
			topic.mu.Lock()
			delete(topic.Clients, w)
			topic.mu.Unlock()
			removeTopicIfEmpty(topicName)

			log.Printf("Client disconnected from topic %s after %ds", topicName, duration)
			return // Disconnect client after timeout
		case <-r.Context().Done(): // Client manually disconnected
			topic.mu.Lock()
			delete(topic.Clients, w)
			topic.mu.Unlock()
			removeTopicIfEmpty(topicName)
			return
		}
	}
}
