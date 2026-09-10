package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Client is the shared MongoDB client used across the application.
var Client *mongo.Client

// Connect initializes the MongoDB client using MONGO_URI (falling back to a
// local default) and verifies connectivity with a ping.
func Connect() *mongo.Client {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("failed to ping mongo: %v", err)
	}

	Client = client
	return client
}

// Database returns a handle to the named database from the shared client.
func Database(name string) *mongo.Database {
	return Client.Database(name)
}
