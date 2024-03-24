package model

import (
	"log"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Userquake struct {
	Code             int
	ObjectID         string
	StartTime        string
	ShortTime        string
	EndTime          string
	AreaByConfidence []AreaByConfidence
}

type AreaByConfidence struct {
	Confidence string
	Areas      []string
}

type UserquakeRecord struct {
	ID              primitive.ObjectID        `bson:"_id"`
	StartedAt       string                    `bson:"started_at"`
	UpdatedAt       string                    `bson:"updated_at"`
	AreaConfidences map[string]AreaConfidence `bson:"area_confidences"`
}

type AreaConfidence struct {
	Confidence float64 `bson:"confidence"`
}

type SortableArea struct {
	Area       string
	Confidence float64
}

func ToUserquake(data primitive.M) (*Userquake, error) {
	log.Printf("ToUserquake: %v", data)

	var uq UserquakeRecord
	bytes, _ := bson.Marshal(data)
	bson.Unmarshal(bytes, &uq)

	return &Userquake{
		Code:             9611,
		ObjectID:         uq.ID.Hex(),
		StartTime:        formatS(uq.StartedAt),
		ShortTime:        formatShort(uq.StartedAt),
		EndTime:          formatTS(uq.UpdatedAt),
		AreaByConfidence: toAreaByConfidence(uq.AreaConfidences),
	}, nil
}

func toAreaByConfidence(ac map[string]AreaConfidence) []AreaByConfidence {
	// 正規化
	max := 0.125
	for _, areaConfidence := range ac {
		if areaConfidence.Confidence > max {
			max = areaConfidence.Confidence
		}
	}

	factor := 1.0 / max
	for area, areaConfidence := range ac {
		ac[area] = AreaConfidence{
			Confidence: areaConfidence.Confidence * factor,
		}
	}

	// 信頼度が高い順に
	var areas []SortableArea
	for area, areaConfidence := range ac {
		areas = append(areas, SortableArea{
			Area:       area,
			Confidence: areaConfidence.Confidence,
		})
	}
	sort.SliceStable(areas, func(i, j int) bool { return areas[i].Confidence > areas[j].Confidence })

	// マップ化
	var abcs []AreaByConfidence
	for _, area := range areas {
		label := confidenceLabel(area.Confidence)
		if label == "F" {
			continue
		}

		if len(abcs) == 0 || abcs[len(abcs)-1].Confidence != label {
			abcs = append(abcs, AreaByConfidence{
				Confidence: label,
				Areas:      []string{},
			})
		}
		abcs[len(abcs)-1].Areas = append(abcs[len(abcs)-1].Areas, area.Area)
	}

	// 地域コード順
	for i := range abcs {
		sort.Strings(abcs[i].Areas)

		var areas []string
		for _, area := range abcs[i].Areas {
			areas = append(areas, convertArea(area))
		}
		abcs[i].Areas = areas
	}

	return abcs
}

func confidenceLabel(confidence float64) string {
	if confidence >= 0.8 {
		return "A"
	} else if confidence >= 0.6 {
		return "B"
	} else if confidence >= 0.4 {
		return "C"
	} else if confidence >= 0.2 {
		return "D"
	} else if confidence >= 0 {
		return "E"
	}
	return "F"
}

var areaMap = map[string]string{
	"900": "地域未設定",
	"901": "地域不明",
	"905": "日本以外",
	"10":  "北海道 石狩",
	"15":  "北海道 渡島",
	"20":  "北海道 檜山",
	"25":  "北海道 後志",
	"30":  "北海道 空知",
	"35":  "北海道 上川",
	"40":  "北海道 留萌",
	"45":  "北海道 宗谷",
	"50":  "北海道 網走",
	"55":  "北海道 胆振",
	"60":  "北海道 日高",
	"65":  "北海道 十勝",
	"70":  "北海道 釧路",
	"75":  "北海道 根室",
	"010": "北海道 石狩",
	"015": "北海道 渡島",
	"020": "北海道 檜山",
	"025": "北海道 後志",
	"030": "北海道 空知",
	"035": "北海道 上川",
	"040": "北海道 留萌",
	"045": "北海道 宗谷",
	"050": "北海道 網走",
	"055": "北海道 胆振",
	"060": "北海道 日高",
	"065": "北海道 十勝",
	"070": "北海道 釧路",
	"075": "北海道 根室",
	"100": "青森津軽",
	"105": "青森三八上北",
	"106": "青森下北",
	"110": "岩手沿岸北部",
	"111": "岩手沿岸南部",
	"115": "岩手内陸",
	"120": "宮城北部",
	"125": "宮城南部",
	"130": "秋田沿岸",
	"135": "秋田内陸",
	"140": "山形庄内",
	"141": "山形最上",
	"142": "山形村山",
	"143": "山形置賜",
	"150": "福島中通り",
	"151": "福島浜通り",
	"152": "福島会津",
	"200": "茨城北部",
	"205": "茨城南部",
	"210": "栃木北部",
	"215": "栃木南部",
	"220": "群馬北部",
	"225": "群馬南部",
	"230": "埼玉北部",
	"231": "埼玉南部",
	"232": "埼玉秩父",
	"240": "千葉北東部",
	"241": "千葉北西部",
	"242": "千葉南部",
	"250": "東京",
	"255": "伊豆諸島北部",
	"260": "伊豆諸島南部",
	"265": "小笠原",
	"270": "神奈川東部",
	"275": "神奈川西部",
	"300": "新潟上越",
	"301": "新潟中越",
	"302": "新潟下越",
	"305": "新潟佐渡",
	"310": "富山東部",
	"315": "富山西部",
	"320": "石川能登",
	"325": "石川加賀",
	"330": "福井嶺北",
	"335": "福井嶺南",
	"340": "山梨東部",
	"345": "山梨中・西部",
	"350": "長野北部",
	"351": "長野中部",
	"355": "長野南部",
	"400": "岐阜飛騨",
	"405": "岐阜美濃",
	"410": "静岡伊豆",
	"411": "静岡東部",
	"415": "静岡中部",
	"416": "静岡西部",
	"420": "愛知東部",
	"425": "愛知西部",
	"430": "三重北中部",
	"435": "三重南部",
	"440": "滋賀北部",
	"445": "滋賀南部",
	"450": "京都北部",
	"455": "京都南部",
	"460": "大阪北部",
	"465": "大阪南部",
	"470": "兵庫北部",
	"475": "兵庫南部",
	"480": "奈良",
	"490": "和歌山北部",
	"495": "和歌山南部",
	"500": "鳥取東部",
	"505": "鳥取中・西部",
	"510": "島根東部",
	"515": "島根西部",
	"514": "島根隠岐",
	"520": "岡山北部",
	"525": "岡山南部",
	"530": "広島北部",
	"535": "広島南部",
	"540": "山口北部",
	"545": "山口中・東部",
	"541": "山口西部",
	"550": "徳島北部",
	"555": "徳島南部",
	"560": "香川",
	"570": "愛媛東予",
	"575": "愛媛中予",
	"576": "愛媛南予",
	"580": "高知東部",
	"581": "高知中部",
	"582": "高知西部",
	"600": "福岡福岡",
	"601": "福岡北九州",
	"602": "福岡筑豊",
	"605": "福岡筑後",
	"610": "佐賀北部",
	"615": "佐賀南部",
	"620": "長崎北部",
	"625": "長崎南部",
	"630": "長崎壱岐・対馬",
	"635": "長崎五島",
	"640": "熊本阿蘇",
	"641": "熊本熊本",
	"645": "熊本球磨",
	"646": "熊本天草・芦北",
	"650": "大分北部",
	"651": "大分中部",
	"655": "大分西部",
	"656": "大分南部",
	"660": "宮崎北部平野部",
	"661": "宮崎北部山沿い",
	"665": "宮崎南部平野部",
	"666": "宮崎南部山沿い",
	"670": "鹿児島薩摩",
	"675": "鹿児島大隅",
	"680": "種子島・屋久島",
	"685": "鹿児島奄美",
	"700": "沖縄本島北部",
	"701": "沖縄本島中南部",
	"702": "沖縄久米島",
	"705": "沖縄八重山",
	"706": "沖縄宮古島",
	"710": "沖縄大東島",
}

func convertArea(code string) string {
	if area, ok := areaMap[code]; ok {
		return area
	}
	return code
}

func formatS(t string) string {
	s, err := time.Parse("2006/01/02 15:04:05.999", t)
	if err != nil {
		return "不明"
	}
	return s.Format("01月02日15時04分05秒")
}

func formatTS(t string) string {
	s, err := time.Parse("2006/01/02 15:04:05.000", t)
	if err != nil {
		return "不明"
	}
	return s.Format("15時04分05秒")
}
