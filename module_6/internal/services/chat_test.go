package services

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupChatTestDB(t *testing.T) *mongo.Database {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Skip("MongoDB not available")
	}
	db := client.Database("test_chat")
	t.Cleanup(func() {
		db.Drop(context.Background())
		client.Disconnect(context.Background())
	})
	return db
}

func TestSaveMessage(t *testing.T) {
	db := setupChatTestDB(t)
	service := NewChatService(db)

	err := service.SaveMessage("user1", "Hello")
	if err != nil {
		t.Fatalf("SaveMessage failed: %v", err)
	}
}

func TestGetHistory(t *testing.T) {
	db := setupChatTestDB(t)
	service := NewChatService(db)

	service.SaveMessage("user1", "Message 1")
	service.SaveMessage("user2", "Message 2")
	service.SaveMessage("user3", "Message 3")

	messages, err := service.GetHistory(10)
	if err != nil {
		t.Fatalf("GetHistory failed: %v", err)
	}

	if len(messages) != 3 {
		t.Fatalf("Expected 3 messages, got %d", len(messages))
	}

	messages, err = service.GetHistory(2)
	if err != nil {
		t.Fatalf("GetHistory failed: %v", err)
	}

	if len(messages) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(messages))
	}
}
