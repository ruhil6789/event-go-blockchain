package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)



type LiquidityModel struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`            // MongoDB Object ID
	SkyContractAddress string             `bson:"skyContractAddress"`        // Sky Contract Address (string)
	Amount            string             `bson:"amount"`                    // Amount (string)
	OwnerAddress      string             `bson:"ownerAddress"`              // Owner Address (string)
	Timestamp         int64              `bson:"timestamp"`                 // Timestamp (int64)
	TransactionHash   string             `bson:"transactionHash"`           // Transaction Hash (string)
	Event             string             `bson:"event"`                     // Event (string, enum)
	BlockNumber       *string            `bson:"blockNumber,omitempty"`     // Block Number (string, nullable)
	CreatedAt         time.Time          `bson:"createdAt,omitempty"`       // Timestamp of creation
	UpdatedAt         time.Time          `bson:"updatedAt,omitempty"`       // Timestamp of last update
}