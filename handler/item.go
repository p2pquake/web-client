package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/p2pquake/web-client/renderer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Service) ItemHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ResponseError(w, http.StatusNotFound, "Not found")
		return
	}
	filters := bson.M{"_id": oid}
	opts := options.FindOneOptions{}

	result := s.Whole.FindOne(context.TODO(), filters, &opts)
	if result.Err() != nil {
		log.Printf("Find error: %v\n", err)
		ResponseError(w, http.StatusNotFound, "Not found")
		return
	}

	var item bson.M
	result.Decode(&item)

	html, err := renderer.RenderItem(item)
	if err != nil {
		log.Printf("Render error: %v\n", err)
		ResponseError(w, http.StatusNotFound, "Not found")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	// w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(html))
}
