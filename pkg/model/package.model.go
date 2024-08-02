package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Package struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"` // MongoDB Object ID
	Amount          int                `bson:"amount"`
	event           string             `bson:"event"`
	added           bool               `bson:"added"`
	transactionHash string             `bson:"added"`
}
