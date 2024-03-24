package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EEW struct {
	Code       int
	ObjectID   string
	Serial     string
	IssueTime  string
	ShortTime  string
	Cancelled  bool
	Hypocenter string
	Areas      []string
}

type EEWRecord struct {
	ID         primitive.ObjectID `bson:"_id"`
	Earthquake struct {
		Hypocenter struct {
			Name string `bson:"name"`
		} `bson:"hypocenter"`
	} `bson:"earthquake"`
	Issue struct {
		Time   string `bson:"time"`
		Serial string `bson:"serial"`
	} `bson:"issue"`
	Cancelled bool      `bson:"cancelled"`
	Areas     []EEWArea `bson:"areas"`
}

type EEWArea struct {
	Pref string `bson:"pref"`
	Name string `bson:"name"`
}

func ToEEW(data primitive.M) (*EEW, error) {
	var eew EEWRecord
	bytes, _ := bson.Marshal(data)
	bson.Unmarshal(bytes, &eew)

	return &EEW{
		ObjectID:   eew.ID.Hex(),
		Code:       556,
		Serial:     eew.Issue.Serial,
		IssueTime:  formatS(eew.Issue.Time),
		ShortTime:  formatShort(eew.Issue.Time),
		Cancelled:  eew.Cancelled,
		Hypocenter: eew.Earthquake.Hypocenter.Name,
		Areas:      toAreas(eew.Areas),
	}, nil
}

func toAreas(areas []EEWArea) []string {
	byPref := make(map[string]bool)
	var result []string
	for _, area := range areas {
		if _, ok := byPref[area.Pref]; !ok {
			byPref[area.Pref] = true
			result = append(result, area.Pref)
		}
	}
	return result
}
