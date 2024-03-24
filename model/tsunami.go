package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tsunami struct {
	Code        int
	ObjectID    string
	Time        string
	IssueTime   string
	ShortTime   string
	Cancelled   bool
	MaxGrade    string
	AreaByGrade []AreaByGrade
}

type AreaByGrade struct {
	Grade string
	Areas []ForecastArea
}

type ForecastArea struct {
	Name        string
	Immediate   bool
	ArrivalTime string
	MaxHeight   string
}

type TsunamiRecord struct {
	ID        primitive.ObjectID `bson:"_id"`
	Time      string             `bson:"time"`
	Cancelled bool               `bson:"cancelled"`
	Issue     Issue              `bson:"issue"`
	Areas     []Area             `bson:"areas"`
}

type Issue struct {
	Time string `bson:"time"`
	Type string `bson:"type"`
}

type Area struct {
	Name        string      `bson:"name"`
	Grade       string      `bson:"grade"`
	Immediate   bool        `bson:"immediate"`
	FirstHeight FirstHeight `bson:"firstHeight"`
	MaxHeight   MaxHeight   `bson:"maxHeight"`
}

type FirstHeight struct {
	ArrivalTime string `bson:"arrivalTime"`
	Condition   string `bson:"condition"`
}

type MaxHeight struct {
	Description string  `bson:"description"`
	Value       float64 `bson:"value"`
}

func ToTsunami(data primitive.M) (*Tsunami, error) {
	var t TsunamiRecord
	bytes, _ := bson.Marshal(data)
	bson.Unmarshal(bytes, &t)

	areaByGrade, maxGrade := toAreaByGrade(t.Areas)

	return &Tsunami{
		Code:        552,
		ObjectID:    t.ID.Hex(),
		Time:        formatS(t.Time),
		IssueTime:   formatS(t.Issue.Time),
		ShortTime:   formatShort(t.Issue.Time),
		Cancelled:   t.Cancelled,
		MaxGrade:    maxGrade,
		AreaByGrade: areaByGrade,
	}, nil
}

func toAreaByGrade(areas []Area) ([]AreaByGrade, string) {
	// 信頼度が高い順に
	gradeEnum := []string{"MajorWarning", "Warning", "Watch", "Unknown"}
	var grades = map[string][]ForecastArea{
		"MajorWarning": {},
		"Warning":      {},
		"Watch":        {},
		"Unknown":      {},
	}
	for _, area := range areas {
		grades[area.Grade] = append(grades[area.Grade], ForecastArea{
			Name:        area.Name,
			Immediate:   area.Immediate,
			ArrivalTime: formatArrivalTime(area.FirstHeight),
			MaxHeight:   formatMaxHeight(area.MaxHeight),
		})
	}

	var result []AreaByGrade
	var max = ""
	for _, grade := range gradeEnum {
		if areas, ok := grades[grade]; ok && len(areas) > 0 {
			result = append(result, AreaByGrade{
				Grade: grade,
				Areas: areas,
			})

			if max == "" {
				max = grade
			}
		}
	}

	return result, max
}

func formatArrivalTime(firstHeight FirstHeight) string {
	if firstHeight.ArrivalTime != "" {
		return formatM(firstHeight.ArrivalTime)
	}

	if firstHeight.Condition == "ただちに津波来襲と推測" {
		return "ただちに来襲"
	}

	if firstHeight.Condition == "津波到達中と推測" {
		return "到達中と推測"
	}

	if firstHeight.Condition == "第１波の到達を確認" {
		return "すでに到達"
	}

	return firstHeight.Condition
}

func formatMaxHeight(maxHeight MaxHeight) string {
	return maxHeight.Description
}

func formatM(t string) string {
	s, err := time.Parse("2006/01/02 15:04:05.999", t)
	if err != nil {
		return "不明"
	}
	return s.Format("15時04分")
}
