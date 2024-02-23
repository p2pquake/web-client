package handler

import (
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	Client *mongo.Client
	Whole  *mongo.Collection
}

func ResponseError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Header().Set("Content-type", "text/plain")
	w.Write([]byte(message))
}
