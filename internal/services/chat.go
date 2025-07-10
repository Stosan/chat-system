package services

import (
	"chatsystem/internal/models"
	"encoding/json"
	"log"
	"net/rpc"

	"gorm.io/gorm"
)

type ChatService struct {
	db *gorm.DB
}

func NewChatService(db *gorm.DB) *ChatService {
	return &ChatService{
		db: db,
	}
}

// Get chat between a user and creator
func (s *ChatService) SaveMessage(msg models.Message) error {
	return nil

}

// PersistMessage sends message to persistence service via RPC
func PersistMessage(msg models.Message) {
	client, err := rpc.DialHTTP("tcp", "localhost:8081")
	if err != nil {
		log.Printf("Error connecting to persistence service: %v", err)
		return
	}
	defer client.Close()

	msgJSON, _ := json.Marshal(msg)
	persistMsg := models.PersistMessage{
		Receiver: msg.Receiver,
		Message:  string(msgJSON),
	}

	var result string
	err = client.Call("PersistenceService.SaveMessage", persistMsg, &result)
	if err != nil {
		log.Printf("Error persisting message: %v", err)
	}
}
