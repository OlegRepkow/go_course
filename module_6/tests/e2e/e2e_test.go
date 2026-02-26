package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestE2EServerRunning(t *testing.T) {
	if os.Getenv("E2E_TEST") != "true" {
		t.Skip("Skipping E2E test. Set E2E_TEST=true to run")
	}

	baseURL := "http://localhost:8080"
	client := &http.Client{Timeout: 5 * time.Second}

	username := fmt.Sprintf("e2euser_%d", time.Now().Unix())

	t.Run("SignUp", func(t *testing.T) {
		body := map[string]string{"username": username, "password": "testpass"}
		jsonBody, _ := json.Marshal(body)

		resp, err := client.Post(baseURL+"/auth/sign-up", "application/json", bytes.NewReader(jsonBody))
		if err != nil {
			t.Fatalf("SignUp request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			t.Fatalf("Expected 200, got %d", resp.StatusCode)
		}
	})

	var token string
	t.Run("SignIn", func(t *testing.T) {
		body := map[string]string{"username": username, "password": "testpass"}
		jsonBody, _ := json.Marshal(body)

		resp, err := client.Post(baseURL+"/auth/sign-in", "application/json", bytes.NewReader(jsonBody))
		if err != nil {
			t.Fatalf("SignIn request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			t.Fatalf("Expected 200, got %d", resp.StatusCode)
		}

		var result map[string]string
		json.NewDecoder(resp.Body).Decode(&result)
		token = result["token"]
		if token == "" {
			t.Fatal("No token received")
		}
	})

	t.Run("SendMessage", func(t *testing.T) {
		body := map[string]string{"text": "E2E test message"}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", baseURL+"/channel/send", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Send request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			t.Fatalf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GetHistory", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/channel/history", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("History request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			t.Fatalf("Expected 200, got %d", resp.StatusCode)
		}
	})
}

func TestE2ECleanup(t *testing.T) {
	if os.Getenv("E2E_TEST") != "true" {
		t.Skip("Skipping E2E cleanup")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Logf("Cleanup: MongoDB connection failed: %v", err)
		return
	}
	defer client.Disconnect(context.Background())

	db := client.Database("chat")
	db.Collection("users").Drop(context.Background())
	db.Collection("messages").Drop(context.Background())
	t.Log("E2E cleanup completed")
}
