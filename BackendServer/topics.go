package BackendServer

import (
	"log"
	"net/http"
	"sync"
)

// Topic structure and global topic map
type Topic struct {
	Name     string
	Messages []Message
	Clients  map[http.ResponseWriter]chan Message
	mu       sync.RWMutex // Mutex for each topic to handle concurrency
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
	if exists {
		return topic
	}

	topicMu.Lock()
	defer topicMu.Unlock()
	newTopic := &Topic{
		Name:    name,
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
			delete(topics, topicName)
			log.Printf("Topic %s removed due to inactivity", topicName)
		}
	}
}
