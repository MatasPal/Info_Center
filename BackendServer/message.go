package BackendServer

import (
	"sync/atomic"
)

// Message structure
type Message struct {
	ID   uint64
	Data string
}

var messageID uint64 // Atomic counter for Message ID

// Generates the next message ID
func getNextMessageID() uint64 {
	return atomic.AddUint64(&messageID, 1)
}
