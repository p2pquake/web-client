package handler

import (
	"context"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/p2pquake/web-client/renderer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Service) IndexHandler(w http.ResponseWriter, r *http.Request) {
	threeDaysAgo := time.Now().Add(time.Hour * -72).Format("2006/01/02 15:04:05")

	// 地震情報・津波予報・緊急地震速報（警報）
	jmaItems, err := s.findJmas(threeDaysAgo)
	if err != nil {
		log.Printf("Render error: %v\n", err)
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 地震感知情報
	userquakeItems, err := s.findUserquakes(threeDaysAgo)
	if err != nil {
		log.Printf("Render error: %v\n", err)
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 並び替え
	var items []bson.M
	items = append(items, jmaItems...)
	items = append(items, userquakeItems...)
	sort.SliceStable(items, func(i, j int) bool {
		return items[i]["time"].(string) > items[j]["time"].(string)
	})

	html, err := renderer.RenderIndex(items)
	if err != nil {
		log.Printf("Render error: %v\n", err)
		ResponseError(w, http.StatusNotFound, "Not found")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	// w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(html))
}

// 地震情報・津波予報・緊急地震速報（警報）
func (s *Service) findJmas(time string) ([]bson.M, error) {
	opts := options.FindOptions{Sort: bson.D{{"$natural", -1}}}
	cursor, err := s.Whole.Find(
		context.TODO(),
		bson.M{
			"code": bson.M{"$in": bson.A{551, 552, 556}},
			"time": bson.M{"$gte": time},
		}, &opts)
	if err != nil {
		return nil, err
	}

	var items []bson.M
	cursor.All(context.TODO(), &items)

	return items, nil
}

// 地震感知情報
func (s *Service) findUserquakes(time string) ([]bson.M, error) {
	opts := options.FindOptions{Sort: bson.D{{"$natural", -1}}}
	cursor, err := s.Whole.Find(
		context.TODO(),
		bson.M{
			"code":       9611,
			"confidence": bson.M{"$gt": 0.9},
			"time":       bson.M{"$gte": time},
		}, &opts)
	if err != nil {
		return nil, err
	}

	var items []bson.M
	cursor.All(context.TODO(), &items)

	// グループ化して除去する必要がある
	var uniqueItems []bson.M
	startedAt := make(map[string]bool)
	for _, item := range items {
		if _, ok := startedAt[item["started_at"].(string)]; !ok {
			startedAt[item["started_at"].(string)] = true
			item["time"] = item["started_at"]
			uniqueItems = append(uniqueItems, item)
		}
	}

	return uniqueItems, nil
}
