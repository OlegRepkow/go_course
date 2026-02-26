package services

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) *mongo.Database {
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

func TestSignUp(t *testing.T) {
	db := setupTestDB(t)
	service := NewAuthService(db, "secret")

	err := service.SignUp("testuser", "password123")
	if err != nil {
		t.Fatalf("SignUp failed: %v", err)
	}
}

func TestSignIn(t *testing.T) {
	db := setupTestDB(t)
	service := NewAuthService(db, "secret")

	service.SignUp("testuser", "password123")

	token, err := service.SignIn("testuser", "password123")
	if err != nil {
		t.Fatalf("SignIn failed: %v", err)
	}
	if token == "" {
		t.Fatal("Expected token")
	}

	_, err = service.SignIn("testuser", "wrongpass")
	if err == nil {
		t.Fatal("Expected invalid credentials error")
	}

	_, err = service.SignIn("nonexistent", "password123")
	if err == nil {
		t.Fatal("Expected invalid credentials error")
	}
}
