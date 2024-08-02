package database

import (
	"context"
	// "fmt"
	// "log"
	// "time"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseConnection struct {
	Client *mongo.Client
	Err    error
}

func ConnectMongoDb() *mongo.Client {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/db")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil
	}

	// Check the connection
	if err = client.Ping(context.TODO(), nil); err != nil {
		return nil
	}

	logrus.Info("MongoClient connected")

	return client
}

// Client instance
var DB *mongo.Client = ConnectMongoDb()

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("event-sky").Collection(collectionName)
	return collection
}
