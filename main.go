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

	whole := client.Database(mongodbDatabase).Collection(mongodbCollection)
	jma := client.Database(mongodbDatabase).Collection("jma")
	service := handler.Service{Client: client, Whole: whole, Jma: jma}

	http.HandleFunc("GET /", service.IndexHandler)
	http.HandleFunc("GET /{id}", service.ItemHandler)
	http.Handle("GET /static/", oneDayCache(http.StripPrefix("/static/", http.FileServer(http.Dir("static")))))

	http.ListenAndServe(":8080", nil)
}

func oneDayCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=86400")
		next.ServeHTTP(w, r)
	})
}
