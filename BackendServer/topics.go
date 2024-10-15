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
	mu       sync.RWMutex
}

var topics = make(map[string]*Topic)

// Retrieves or creates a topic
func getTopic(name string) *Topic {
	mu.Lock()
	defer mu.Unlock()

	if topic, exists := topics[name]; exists {
		return topic
	}

	newTopic := &Topic{
		Name:    name,
		Clients: make(map[http.ResponseWriter]chan Message),
	}
	topics[name] = newTopic
	return newTopic
}

// Removes a topic if no clients are left
func removeTopicIfEmpty(topicName string) {
	mu.Lock()
	defer mu.Unlock()

	topic, exists := topics[topicName]
	if exists && len(topic.Clients) == 0 {
		delete(topics, topicName)
		log.Printf("Topic %s removed due to inactivity", topicName)
	}
}
