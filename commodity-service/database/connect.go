package database

import (
	"commodity-service/config"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client holds the MongoDB client instance.
var Client *mongo.Client

// ConnectDB establishes a connection to MongoDB using config.Cfg.MongoDBURI.
func ConnectDB() (*mongo.Client, error) { // No arguments here
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Cfg.MongoDBURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Successfully connected to MongoDB!")
	Client = client
	return client, nil
}

// GetCollection returns a handle to a MongoDB collection.
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database(config.Cfg.DatabaseName).Collection(collectionName)
}
