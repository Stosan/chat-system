package websocket

import "chatsystem/internal/services"

// ProcessChatMessages processes chat messages from channel
func (s *Hub) ProcessChatMessages() {
	for msg := range s.chatChan {
		if receiver, exists := s.GetClient(msg.Receiver); exists {
			receiver.Conn.WriteJSON(msg)
		}
	}
}

// ProcessPersistMessages processes persistence messages from channel
func (s *Hub) ProcessPersistMessages() {
	for msg := range s.persistChan {
		go services.PersistMessage(msg)
	}
}
