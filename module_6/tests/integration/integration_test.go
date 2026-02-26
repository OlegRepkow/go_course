package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"module_6/cmd/server/handlers"
	"module_6/cmd/server/middlewares"
	"module_6/internal/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupIntegrationTest(t *testing.T) (*fiber.App, *mongo.Database) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Skip("MongoDB not available")
	}
	db := client.Database("test_integration")
	t.Cleanup(func() {
		db.Drop(context.Background())
		client.Disconnect(context.Background())
	})

	secret := "testsecret"
	authService := services.NewAuthService(db, secret)
	chatService := services.NewChatService(db)

	authHandler := handlers.NewAuthHandler(authService)
	chatHandler := handlers.NewChatHandler(chatService)

	app := fiber.New()
	app.Post("/auth/sign-up", authHandler.SignUp)
	app.Post("/auth/sign-in", authHandler.SignIn)

	protected := app.Group("", middlewares.Auth(secret))
	protected.Get("/channel/history", chatHandler.History)
	protected.Post("/channel/send", chatHandler.Send)

	return app, db
}

func TestFullAuthFlow(t *testing.T) {
	app, _ := setupIntegrationTest(t)

	signUpBody := map[string]string{"username": "testuser", "password": "pass123"}
	jsonBody, _ := json.Marshal(signUpBody)

	req := httptest.NewRequest("POST", "/auth/sign-up", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Fatalf("SignUp failed: %d", resp.StatusCode)
	}

	signInBody := map[string]string{"username": "testuser", "password": "pass123"}
	jsonBody, _ = json.Marshal(signInBody)

	req = httptest.NewRequest("POST", "/auth/sign-in", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)

	if resp.StatusCode != 200 {
		t.Fatalf("SignIn failed: %d", resp.StatusCode)
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	if result["token"] == "" {
		t.Fatal("No token returned")
	}
}

func TestFullChatFlow(t *testing.T) {
	app, _ := setupIntegrationTest(t)

	signUpBody := map[string]string{"username": "chatuser", "password": "pass123"}
	jsonBody, _ := json.Marshal(signUpBody)
	req := httptest.NewRequest("POST", "/auth/sign-up", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	app.Test(req)

	signInBody := map[string]string{"username": "chatuser", "password": "pass123"}
	jsonBody, _ = json.Marshal(signInBody)
	req = httptest.NewRequest("POST", "/auth/sign-in", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	var authResult map[string]string
	json.NewDecoder(resp.Body).Decode(&authResult)
	token := authResult["token"]

	sendBody := map[string]string{"text": "Hello integration test"}
	jsonBody, _ = json.Marshal(sendBody)
	req = httptest.NewRequest("POST", "/channel/send", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ = app.Test(req)

	if resp.StatusCode != 200 {
		t.Fatalf("Send failed: %d", resp.StatusCode)
	}

	req = httptest.NewRequest("GET", "/channel/history", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ = app.Test(req)

	if resp.StatusCode != 200 {
		t.Fatalf("History failed: %d", resp.StatusCode)
	}
}
