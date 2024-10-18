package BackendServer

import (
	"log"
	"net/http"
	"sync"
)

// Topic structure and global topic map
type Topic struct {
	Messages []Message
	Clients  map[http.ResponseWriter]chan Message
	mu       sync.RWMutex // Mutex for each topic to handle concurrency
	removed  bool         // Flag to indicate if the topic has been removed
}

var (
	topics  = make(map[string]*Topic)
	topicMu sync.RWMutex // Mutex for managing topics map
)

// Retrieves or creates a topic if it doesn’t already exist
func getTopic(name string) *Topic {
	topicMu.RLock()
	topic, exists := topics[name]
	topicMu.RUnlock()

	if exists && topic.removed {
		// If the topic was removed, return nil
		return nil
	}

	if exists {
		return topic
	}

	topicMu.Lock()
	defer topicMu.Unlock()

	// Check again in case it was created during the lock acquisition
	topic, exists = topics[name]
	if exists {
		return topic
	}

	newTopic := &Topic{
		Clients: make(map[http.ResponseWriter]chan Message),
	}
	topics[name] = newTopic
	return newTopic
}

// Removes a topic if no clients are left, ensuring proper locking
func removeTopicIfEmpty(topicName string) {
	topicMu.Lock()
	defer topicMu.Unlock()

	topic, exists := topics[topicName]
	if exists {
		topic.mu.Lock()
		defer topic.mu.Unlock()
		if len(topic.Clients) == 0 {
			topic.removed = true // Mark the topic as removed
			delete(topics, topicName)
			log.Printf("Topic %s removed due to inactivity", topicName)
		}
	}
}
