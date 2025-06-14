package handler

import (
	"context"
	"encoding/json"
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

func (s *Service) TimeseriesHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	log.Printf("TimeseriesHandler called with ID: %s", id)
	
	if id == "" {
		log.Printf("Empty ID received")
		ResponseError(w, http.StatusBadRequest, "Empty ID")
		return
	}
	
	// 指定されたIDに対応するstarted_atを取得
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ObjectID format: %s, error: %v", id, err)
		ResponseError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	// 最初のレコードを取得してstarted_atを確認
	var firstItem bson.M
	err = s.Whole.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&firstItem)
	if err != nil {
		log.Printf("Item not found for ID: %s, error: %v", id, err)
		ResponseError(w, http.StatusNotFound, "Item not found")
		return
	}

	code, ok := firstItem["code"]
	log.Printf("Found item with code: %v (type: %T)", code, code)
	
	// コードの型に応じて比較
	var codeInt int
	switch v := code.(type) {
	case int:
		codeInt = v
	case int32:
		codeInt = int(v)
	case int64:
		codeInt = int(v)
	case float64:
		codeInt = int(v)
	default:
		log.Printf("Unexpected code type: %T, value: %v", code, code)
		ResponseError(w, http.StatusBadRequest, "Invalid code type")
		return
	}
	
	if !ok || codeInt != 9611 {
		log.Printf("Not a userquake event, code: %v (int: %d)", code, codeInt)
		ResponseError(w, http.StatusBadRequest, "Not a userquake event")
		return
	}
	
	log.Printf("Confirmed userquake event with code: %d", codeInt)

	startedAt, ok := firstItem["started_at"].(string)
	if !ok {
		log.Printf("Invalid started_at field: %v", firstItem["started_at"])
		ResponseError(w, http.StatusInternalServerError, "Invalid started_at")
		return
	}
	log.Printf("Found started_at: %s", startedAt)

	// 同じstarted_atを持つ全てのレコードを取得
	opts := options.FindOptions{Sort: bson.D{{"updated_at", 1}}}
	cursor, err := s.Whole.Find(
		context.TODO(),
		bson.M{
			"code":       9611,
			"started_at": startedAt,
		}, &opts)
	if err != nil {
		log.Printf("Database find error: %v", err)
		ResponseError(w, http.StatusInternalServerError, "Database error")
		return
	}

	var items []bson.M
	if err = cursor.All(context.TODO(), &items); err != nil {
		log.Printf("Database cursor error: %v", err)
		ResponseError(w, http.StatusInternalServerError, "Database error")
		return
	}

	log.Printf("Found %d timeseries items for started_at: %s", len(items), startedAt)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
