package renderer

import (
	"github.com/p2pquake/web-client/model"
	"go.mongodb.org/mongo-driver/bson"
)

func RenderItem(m bson.M) (string, error) {
	data, err := model.Convert(m)
	if err != nil {
		return "", err
	}

	return Render("item.html", data)
}
