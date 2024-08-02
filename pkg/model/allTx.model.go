package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AllTxModel struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`           // MongoDB Object ID
	ReceiverAddress  string             `bson:"receiverAddress"`         // Receiver's address (string)
	Amount           string             `bson:"amount"`                  // Amount (string)
	Amt              float64            `bson:"amt"`                     // Amount in numeric format (number)
	Timestamp        *int64             `bson:"timestamp,omitempty"`     // Timestamp (number), nullable
	TransactionHash  string             `bson:"transactionHash"`         // Transaction hash (string)
	Event            string             `bson:"event"`                   // Event type (string, enum)
	BlockNumber      string             `bson:"blockNumber,omitempty"`   // Block number (string), nullable
}