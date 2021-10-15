package handler

import (
	"context"
	_ "crypto/sha1"
	"fmt"
	"net/http"
	_ "strings"
	"time"

	"github.com/labstack/echo"
	_ "github.com/ponyo877/news-app-backend/handler/imagectl"
	"github.com/ponyo877/news-app-backend/handler/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	orgmongo "go.mongodb.org/mongo-driver/mongo"
)

// func deviceID2sha1(imageUrl string) uint32 {
// 	h := sha1.New()
// 	h.Write([]byte(imageUrl))
// 	bs := h.Sum()
// 	return bs
// }

// https://stackoverflow.com/questions/37407476/how-to-push-in-golang-for-nested-array
// curl -XPOST -F "articleID=<AID>" -F "massage=<MSG>" -F "user=<USR>" -F "dID=<DID>" https://gitouhon-juku-k8s2.ga/comment/put
func PutComments() echo.HandlerFunc {
	return func(c echo.Context) error {
		qarticleIDTmp := c.FormValue("articleID")
		qmassage := c.FormValue("massage")
		qdeviceHash := c.FormValue("devicehash")
		// qdeviceHash := fmt.Sprint(imagectl.HashID(qdeviceID))
		fmt.Println(qarticleIDTmp)
		fmt.Println(qmassage)
		// fmt.Println(qdeviceID)
		qarticleID, err := primitive.ObjectIDFromHex(qarticleIDTmp)
		checkError(err)

		olayput := "2006-01-02 15:04:05.000"
		t := time.Now()
		qpostDate := t.Format(olayput)

		ctx := context.Background()
		client := mongo.OpenMongo()
		err = client.Connect(ctx)
		defer client.Disconnect(ctx)
		checkError(err)

		var usr_info map[string]interface{}
		usr_col := client.Database("newsdb").Collection("user_col")
		usr_filter := bson.M{
			"deviceHash": bson.M{"$eq": qdeviceHash},
		}
		err = usr_col.FindOne(ctx, usr_filter).Decode(&usr_info)
		if err == orgmongo.ErrNoDocuments {
			return c.JSON(http.StatusOK, map[string]string{"Status": "NG"})
		}
		checkError(err)

		art_col := client.Database("newsdb").Collection("article_col")
		newComment := map[string]interface{}{
			"username":   usr_info["name"],
			"avatar":     usr_info["avatar"],
			"deviceHash": usr_info["deviceHash"],
			"massage":    qmassage,
			"postDate":   qpostDate,
		}
		art_filter := bson.M{
			"_id": bson.M{"$eq": qarticleID},
		}
		art_change := bson.M{
			"$push": bson.M{
				"comments": newComment,
			},
		}
		_, err = art_col.UpdateOne(ctx, art_filter, art_change)
		checkError(err)
		return c.JSON(http.StatusOK, map[string]string{"Status": "Ok"})
	}
}
