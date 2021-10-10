package handler

import (
	"context"
	_"fmt"
	"net/http"
	_"strings"
	"./mongo"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetComments() echo.HandlerFunc {
	return func(c echo.Context) error {
		qarticleIDTmp := c.QueryParam("articleID")
		qarticleID, err := primitive.ObjectIDFromHex(qarticleIDTmp)
		checkError(err)

		ctx := context.Background()
		client := mongo.OpenMongo()
		err = client.Connect(ctx)
		defer client.Disconnect(ctx)
		checkError(err)

		col := client.Database("newsdb").Collection("article_col")

		filter := bson.M{
			"_id": bson.M{"$eq": qarticleID},
		}
		var feedmap map[string]interface{}
		err = col.FindOne(ctx, filter).Decode(&feedmap)
		checkError(err)
		
		comments := feedmap["comments"]
		return c.JSON(http.StatusOK, map[string]interface{}{"data": comments})
	}
}

