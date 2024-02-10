package handler

import (
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	Client *mongo.Client
	Whole  *mongo.Collection
}

func (s *Service) IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}
