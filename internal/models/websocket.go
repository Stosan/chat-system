package models





// PersistMessage represents message for persistence
type PersistMessage struct {
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}

// QueryMessage represents query for historical messages
type QueryMessage struct {
	UserID string `json:"user_id"`
}
