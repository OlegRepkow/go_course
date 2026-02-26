package services

import (
	"context"
	"log"
	"module_6/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatService struct {
	messages *mongo.Collection
}

func NewChatService(db *mongo.Database) *ChatService {
	return &ChatService{messages: db.Collection("messages")}
}

func (s *ChatService) SaveMessage(username, text string) error {
	log.Printf("[INFO] ChatService.SaveMessage: user=%s, text=%s\n", username, text)
	msg := models.Message{
		Username:  username,
		Text:      text,
		Timestamp: time.Now(),
	}
	_, err := s.messages.InsertOne(context.Background(), msg)
	if err != nil {
		log.Printf("[ERROR] ChatService.SaveMessage: DB insert failed: %v\n", err)
		return err
	}
	log.Printf("[INFO] ChatService.SaveMessage: message saved successfully\n")
	return nil
}

func (s *ChatService) GetHistory(limit int) ([]models.Message, error) {
	log.Printf("[INFO] ChatService.GetHistory: fetching last %d messages\n", limit)
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(int64(limit))
	cursor, err := s.messages.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		log.Printf("[ERROR] ChatService.GetHistory: DB query failed: %v\n", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var messages []models.Message
	if err := cursor.All(context.Background(), &messages); err != nil {
		log.Printf("[ERROR] ChatService.GetHistory: cursor decode failed: %v\n", err)
		return nil, err
	}
	log.Printf("[INFO] ChatService.GetHistory: fetched %d messages\n", len(messages))
	return messages, nil
}
