package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type EventBlock struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`    // MongoDB Object ID
	Chain       string             `bson:"chain"`            // Chain (string)
	Address     string             `bson:"address"`          // Address (string)
	BlockNumber *int64             `bson:"blockNumber,omitempty"` // Block Number (number), nullable
	CreatedAt   time.Time          `bson:"createdAt,omitempty"`   // Timestamp of creation
	UpdatedAt   time.Time          `bson:"updatedAt,omitempty"`   // Timestamp of last update
}