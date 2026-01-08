package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Collections struct {
	UserCollection       *mongo.Collection
	CredentialCollection *mongo.Collection
	ApiKeyCollection     *mongo.Collection
	TicketCollection     *mongo.Collection
}

var DbCollections Collections
var dbClient *mongo.Client

func InitCollections() error {
	err := InitDB()
	if err != nil {
		return err // Return error instead of log.Fatal
	}
	DbCollections.UserCollection = GetDBCollection("users")
	DbCollections.CredentialCollection = GetDBCollection("credentials")
	DbCollections.ApiKeyCollection = GetDBCollection("api_keys")
	DbCollections.TicketCollection = GetDBCollection("tickets")
	return nil
}

func CloseDB() {
	if dbClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := dbClient.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		} else {
			log.Println("Disconnected from MongoDB.")
		}
	}
}

func InitDB() error {
	url := os.Getenv("MONGO_URI")
	clientOptions := options.Client().
		ApplyURI(url).
		SetServerSelectionTimeout(30 * time.Second). // Increase timeout
		SetConnectTimeout(10 * time.Second).
		SetSocketTimeout(10 * time.Second).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(30 * time.Second).
		SetDirect(false) // Important for replica sets

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error
	dbClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err // Return error instead of log.Fatal
	}

	// Ping with primary read preference (important for replica sets)
	err = dbClient.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	log.Println("Connected to MongoDB!")

	return nil
}

func GetDBCollection(collectionName string) *mongo.Collection {
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "mydb" // fallback
	}
	return dbClient.Database(dbName).Collection(collectionName)
}

// GetClient returns the MongoDB client for transaction usage
func GetClient() *mongo.Client {
	return dbClient
}

// WithTransaction provides a helper for running transactions
func WithTransaction(client *mongo.Client, ctx context.Context, fn func(sessCtx mongo.SessionContext) error) error {
	if client == nil {
		return mongo.ErrClientDisconnected
	}

	session, err := client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, fn(sessCtx)
	}

	_, err = session.WithTransaction(ctx, callback)
	return err
}
