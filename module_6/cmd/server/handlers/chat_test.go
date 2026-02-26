package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"module_6/internal/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupChatTestDB(t *testing.T) *mongo.Database {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Skip("MongoDB not available")
	}
	db := client.Database("test_chat_handlers")
	t.Cleanup(func() {
		db.Drop(context.Background())
		client.Disconnect(context.Background())
	})
	return db
}

func TestHistoryHandler(t *testing.T) {
	db := setupChatTestDB(t)
	service := services.NewChatService(db)
	handler := NewChatHandler(service)

	service.SaveMessage("user1", "Test message")

	app := fiber.New()
	app.Get("/history", func(c *fiber.Ctx) error {
		c.Locals("username", "user1")
		return handler.History(c)
	})

	req := httptest.NewRequest("GET", "/history", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestSendHandler(t *testing.T) {
	db := setupChatTestDB(t)
	service := services.NewChatService(db)
	handler := NewChatHandler(service)

	app := fiber.New()
	app.Post("/send", func(c *fiber.Ctx) error {
		c.Locals("username", "user1")
		return handler.Send(c)
	})

	body := map[string]string{"text": "Hello world"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/send", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}
}
