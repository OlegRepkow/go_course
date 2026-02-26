package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"module_6/internal/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
)

func setupTestDB(t *testing.T) *mongo.Database {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Skip("MongoDB not available")
	}
	db := client.Database("test_handlers")
	t.Cleanup(func() {
		db.Drop(context.Background())
		client.Disconnect(context.Background())
	})
	return db
}

func TestSignUpHandler(t *testing.T) {
	db := setupTestDB(t)
	service := services.NewAuthService(db, "secret")
	handler := NewAuthHandler(service)

	app := fiber.New()
	app.Post("/sign-up", handler.SignUp)

	body := map[string]string{"username": "testuser", "password": "pass123"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/sign-up", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestSignInHandler(t *testing.T) {
	db := setupTestDB(t)
	service := services.NewAuthService(db, "secret")
	handler := NewAuthHandler(service)

	service.SignUp("testuser", "pass123")

	app := fiber.New()
	app.Post("/sign-in", handler.SignIn)

	body := map[string]string{"username": "testuser", "password": "pass123"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/sign-in", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	if result["token"] == "" {
		t.Fatal("Expected token in response")
	}
}

func TestSignInHandlerInvalidCredentials(t *testing.T) {
	db := setupTestDB(t)
	service := services.NewAuthService(db, "secret")
	handler := NewAuthHandler(service)

	app := fiber.New()
	app.Post("/sign-in", handler.SignIn)

	body := map[string]string{"username": "nonexistent", "password": "wrong"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/sign-in", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != 401 {
		t.Fatalf("Expected 401, got %d", resp.StatusCode)
	}
}
