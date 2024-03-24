package model

import "go.mongodb.org/mongo-driver/bson"

func Convert(data bson.M) (interface{}, error) {
	var result interface{}
	var err error = nil
	switch toInt(data["code"]) {
	case 551:
		result, err = ToEarthquake(data)
	case 552:
		result, err = ToTsunami(data)
	case 556:
		result, err = ToEEW(data)
	case 9611:
		result, err = ToUserquake(data)
	default:
		result = data
	}
	return result, err
}

func toInt(e interface{}) int {
	switch e.(type) {
	case int:
		return e.(int)
	case int32:
		return int(e.(int32))
	case int64:
		return int(e.(int64))
	}
	return 0
}
