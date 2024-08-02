package model

import(
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)


type RewardIncome struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`        // MongoDB Object ID
	Percent         float64            `bson:"percent"`              // Percent (Number)
	Level           int                `bson:"level"`                // Level (Number)
	Timestamp       string             `bson:"timestamp"`            // Timestamp (String)
	Event           string             `bson:"event"`                // Event (String)
	TransactionHash string             `bson:"transactionHash"`      // Transaction Hash (String)
	BlockNumber     string             `bson:"blockNumber,omitempty"`// Block Number (String)
	CreatedAt       time.Time          `bson:"createdAt,omitempty"`  // Timestamp of creation
	UpdatedAt       time.Time          `bson:"updatedAt,omitempty"`  // Timestamp of last update
}