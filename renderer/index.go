package renderer

import (
	"github.com/p2pquake/web-client/model"
	"go.mongodb.org/mongo-driver/bson"
)

func RenderIndex(ms []bson.M) (string, error) {
	data := make([]interface{}, len(ms))
	var err error
	for i, m := range ms {
		data[i], err = model.Convert(m)
		if err != nil {
			return "", err
		}
	}

	return Render("index.html", data)
}
