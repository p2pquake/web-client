package model

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Earthquake struct {
	Raw              string
	Code             int
	ObjectID         string
	MaxScale         string
	IssueTime        string
	IssueType        string
	OccurredTime     string
	ShortTime        string
	Hypocenter       string
	IsEruption       bool
	FreeFormComments []string
	Tsunami          string
	ForeignTsunami   string
	Points           []PointsByPref
	PointsByScale    []PointsByScale
}

type PointsByPref struct {
	Pref   string
	Points []PointsByScale
}

type PointsByScale struct {
	Scale  string
	Points []string
}

func (ps PointsByScale) PointString() string {
	return strings.Join(ps.Points, "、")
}

type EarthquakeRecord struct {
	ID         primitive.ObjectID `bson:"_id"`
	Earthquake struct {
		MaxScale        int        `bson:"maxScale"`
		Time            string     `bson:"time"`
		DomesticTsunami string     `bson:"domesticTsunami"`
		ForeignTsunami  string     `bson:"foreignTsunami"`
		Hypocenter      Hypocenter `bson:"hypocenter"`
	} `bson:"earthquake"`
	Issue struct {
		Type string `bson:"type"`
		Time string `bson:"time"`
	} `bson:"issue"`
	Points   []Point `bson:"points"`
	Comments struct {
		FreeFormComment string `bson:"freeFormComment"`
	} `bson:"comments"`
}

type Point struct {
	Pref   string `bson:"pref"`
	Addr   string `bson:"addr"`
	IsArea bool   `bson:"isArea"`
	Scale  int    `bson:"scale"`
}

type Hypocenter struct {
	Name      string  `bson:"name"`
	Depth     int     `bson:"depth"`
	Latitude  float64 `bson:"latitude"`
	Longitude float64 `bson:"longitude"`
	Magnitude float64 `bson:"magnitude"`
}

func ToEarthquake(data primitive.M) (*Earthquake, error) {
	var eq EarthquakeRecord
	bytes, _ := bson.Marshal(data)
	bson.Unmarshal(bytes, &eq)

	// 「震度5弱以上と推定」の優先度を下げる（震度5弱より低い）
	for i := 0; i < len(eq.Points); i++ {
		if eq.Points[i].Scale == 46 {
			eq.Points[i].Scale = 44
		}
	}

	r := regexp.MustCompile("^((?:余市町|田村市|玉村町|東村山市|武蔵村山市|羽村市|十日町市|上市町|大町市|名古屋中村区|大阪堺市.+?区|下市町|大村市|野々市市|四日市市|廿日市市|大町町|.+?[市区町村]))")
	sort.SliceStable(eq.Points, func(i, j int) bool { return eq.Points[i].Scale > eq.Points[j].Scale })
	for i, v := range eq.Points {
		s := r.FindString(v.Addr)
		if s != "" {
			eq.Points[i] = Point{
				Pref:   v.Pref,
				Addr:   s,
				IsArea: v.IsArea,
				Scale:  v.Scale,
			}
		}
	}

	// 同じ市区町村名の場合、震度の大きいものを残す
	var points []Point
	byCities := make(map[string]Point)
	for _, p := range eq.Points {
		if _, ok := byCities[p.Addr]; !ok {
			byCities[p.Addr] = p
			points = append(points, p)
		}
	}

	// 都道府県ごとにグループ化
	var prefs []string
	byPrefs := make(map[string][]Point)
	for _, p := range points {
		if _, ok := byPrefs[p.Pref]; !ok {
			byPrefs[p.Pref] = make([]Point, 0)
			prefs = append(prefs, p.Pref)
		}
		byPrefs[p.Pref] = append(byPrefs[p.Pref], p)
	}

	var pointsByPref []PointsByPref
	for _, pref := range prefs {
		var scales []int
		byScale := make(map[int][]string)
		for _, p := range byPrefs[pref] {
			if _, ok := byScale[p.Scale]; !ok {
				byScale[p.Scale] = make([]string, 0)
				scales = append(scales, p.Scale)
			}
			byScale[p.Scale] = append(byScale[p.Scale], p.Addr)
		}

		var pointsByScale []PointsByScale
		for _, s := range scales {
			pointsByScale = append(pointsByScale, PointsByScale{
				Scale:  scale(s),
				Points: byScale[s],
			})
		}

		pointsByPref = append(pointsByPref, PointsByPref{
			Pref:   pref,
			Points: pointsByScale,
		})
	}

	// 震度速報に限り、震度の大きい順の情報を提供
	var scales []int
	byScale := make(map[int][]string)
	for _, p := range points {
		if _, ok := byScale[p.Scale]; !ok {
			byScale[p.Scale] = make([]string, 0)
			scales = append(scales, p.Scale)
		}
		byScale[p.Scale] = append(byScale[p.Scale], p.Addr)
	}

	var pointsByScale []PointsByScale
	for _, s := range scales {
		pointsByScale = append(pointsByScale, PointsByScale{
			Scale:  scale(s),
			Points: byScale[s],
		})
	}

	isEruption := strings.Contains(eq.Comments.FreeFormComment, "大規模な噴火")
	freeFormComments := []string{}
	if eq.Comments.FreeFormComment != "" {
		freeFormComments = strings.Split(eq.Comments.FreeFormComment, "\n")
	}

	return &Earthquake{
		Raw:              fmt.Sprintf("%v\n", eq),
		ObjectID:         eq.ID.Hex(),
		Code:             551,
		MaxScale:         scale(eq.Earthquake.MaxScale),
		IssueType:        eq.Issue.Type,
		IssueTime:        format(eq.Issue.Time),
		OccurredTime:     format(eq.Earthquake.Time),
		ShortTime:        formatShort(eq.Earthquake.Time),
		Tsunami:          tsunami(eq.Earthquake.DomesticTsunami),
		ForeignTsunami:   tsunami(eq.Earthquake.ForeignTsunami),
		Hypocenter:       hypocenter(eq.Earthquake.Hypocenter, isEruption),
		IsEruption:       isEruption,
		FreeFormComments: freeFormComments,
		Points:           pointsByPref,
		PointsByScale:    pointsByScale,
	}, nil
}

func scale(s int) string {
	switch s {
	case 10:
		return "1"
	case 20:
		return "2"
	case 30:
		return "3"
	case 40:
		return "4"
	case 45:
		return "5弱"
	case 44:
		return "5弱以上と推定"
	case 46: // 46 -> 44 に書き換えているのでおそらく不要。
		return "5弱以上と推定"
	case 50:
		return "5強"
	case 55:
		return "6弱"
	case 60:
		return "6強"
	case 70:
		return "7"
	}
	return "不明"
}

func format(t string) string {
	s, err := time.Parse("2006/01/02 15:04:05", t)
	if err != nil {
		return "不明"
	}
	return s.Format("01月02日15時04分頃")
}

func formatShort(t string) string {
	s, err := time.Parse("2006/01/02 15:04:05", t)
	if err != nil {
		return "不明"
	}
	return s.Format("01/02 15:04頃")
}

func tsunami(t string) string {
	switch t {
	case "None":
		return "津波の心配なし"
	case "Unknown":
		return "津波有無は不明"
	case "Checking":
		return "津波有無は調査中"
	case "NonEffective":
		return "津波被害の心配なし（若干の海面変動あり）"
	case "Watch":
		return "津波注意報 発表中"
	case "Warning":
		return "津波予報 発表中"
	case "NonEffectiveNearby":
		return "津波被害の心配なし（震源近傍で小さな津波の可能性あり）"
	case "WarningNearby":
		return "震源近傍で津波の可能性あり"
	case "WarningPacific":
		return "太平洋で津波の可能性あり"
	case "WarningPacificWide":
		return "太平洋広域で津波の可能性あり"
	case "WarningIndian":
		return "インド洋で津波の可能性あり"
	case "WarningIndianWide":
		return "インド洋広域で津波の可能性あり"
	case "Potential":
		return "この規模は一般的に津波の可能性あり"
	}
	return "津波有無は不明"
}

func hypocenter(hypocenter Hypocenter, isEruption bool) string {
	if hypocenter.Name == "" {
		return "不明"
	}

	if isEruption {
		return hypocenter.Name
	}

	return fmt.Sprintf("%s (%s) M%.1f", hypocenter.Name, depth(hypocenter.Depth), hypocenter.Magnitude)
}

func depth(depth int) string {
	if depth < 0 {
		return "深さ不明"
	}
	if depth == 0 {
		return "ごく浅い深さ"
	}
	return fmt.Sprintf("深さ%dkm", depth)
}
