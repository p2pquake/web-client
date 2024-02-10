package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Service) ItemHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("ObjectID decode error (%v): %v\n", id, err)
		return
	}
	filters := bson.M{"_id": oid}
	// opts := options.FindOneOptions{Sort: bson.D{{"time", -1}}}
	opts := options.FindOneOptions{}

	result := s.Whole.FindOne(context.TODO(), filters, &opts)

	var item bson.M
	result.Decode(&item)

	fmt.Fprintf(w, "ObjectID: %v\n", oid)
	fmt.Fprintf(w, "Item: %v", item)
}
