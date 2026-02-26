package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Message represents a chat message
type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username"`
	Text      string             `bson:"text" json:"text"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}
