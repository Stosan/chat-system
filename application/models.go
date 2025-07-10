package application

import "time"

// Message represents a chat message
type Message struct {
	Text      string    `json:"text"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

// PersistMessage represents message for persistence
type PersistMessage struct {
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}

// QueryMessage represents query for historical messages
type QueryMessage struct {
	UserID string `json:"user_id"`
}
