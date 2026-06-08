package config

import (
    "context"
    "fmt"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "github.com/joho/godotenv"
)

var Database *mongo.Database

func init() {
    _ = godotenv.Load()

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    mongoURI := os.Getenv("MONGOSTRING")
    if mongoURI == "" {
        fmt.Println("MONGOSTRING environment variable is not set")
        os.Exit(1)
    }

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
    if err != nil {
        fmt.Printf("Failed to connect to MongoDB: %v\n", err)
        os.Exit(1)
    }

    if err = client.Ping(ctx, nil); err != nil {
        fmt.Printf("Unable to ping MongoDB: %v\n", err)
        os.Exit(1)
    }

    Database = client.Database("airdropv2")
    fmt.Println("Successfully connected to MongoDB")
}