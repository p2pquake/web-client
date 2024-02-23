package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/p2pquake/web-client/handler"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongodbUrl := os.Getenv("MONGODB_URL")
	mongodbDatabase := os.Getenv("DATABASE")
	mongodbCollection := os.Getenv("COLLECTION")

	opts := options.Client().ApplyURI(mongodbUrl)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("P2PQuake web client")
	log.Printf("Database %v, Collection: %v\n", mongodbDatabase, mongodbCollection)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatalf("MongoDB connect error: %v", err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database(mongodbDatabase).Collection(mongodbCollection)
	service := handler.Service{Client: client, Whole: collection}

	http.HandleFunc("GET /{id}", service.ItemHandler)

	http.ListenAndServe(":8080", nil)
}
