package renderer

import "go.mongodb.org/mongo-driver/bson"

func RenderItem(data bson.M) (string, error) {
	return Render("item.html", data)
}
