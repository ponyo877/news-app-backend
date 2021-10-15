package handler

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/labstack/echo"
	"github.com/ponyo877/news-app-backend/handler/mongo"
	"github.com/ponyo877/news-app-backend/handler/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	orgmongo "go.mongodb.org/mongo-driver/mongo"
	_ "go.mongodb.org/mongo-driver/mongo/options"
)

func GetRankingMongo() echo.HandlerFunc {
	return func(c echo.Context) error {
		kindStr := c.Param("kind")
		fmt.Println("kindStr: ", kindStr)
		fmt.Println("reflect.TypeOf(kindStr): ", reflect.TypeOf(kindStr))
		ctx := context.Background()
		client := mongo.OpenMongo()
		err := client.Connect(ctx)
		defer client.Disconnect(ctx)
		checkError(err)

		idsRanking := redis.GetIdsRankingTmp(kindStr)
		col := client.Database("newsdb").Collection("article_col")
		// var feed map[string]interface{}
		var feedArray []map[string]interface{}

		var publishedAtTmp string
		ilayout := "2006-01-02T15:04:05-07:00"
		olayput := "2006-01-02 15:04"

		for _, id_count := range idsRanking {
			var feed map[string]interface{}
			_id, err := primitive.ObjectIDFromHex(id_count["id"].(string))
			if err == hex.ErrLength {
				continue
			}
			checkError(err)
			filter := bson.M{
				"_id": bson.M{"$eq": _id},
			}
			err = col.FindOne(ctx, filter).Decode(&feed) // , findOptions)
			if err == orgmongo.ErrNoDocuments {
				continue
			}
			checkError(err)

			// change time format
			publishedAtTmp = feed["publishedAt"].(string)
			tfm, _ := time.Parse(ilayout, publishedAtTmp)
			// feed["publishedAt"] = tfm.Format(olayput)
			feedmap := map[string]interface{}{
				"_id":         feed["_id"],
				"viewcount":   id_count["viewcount"], //strconv.Itoa(id_count["viewcount"]),
				"image":       feed["image"],
				"publishedAt": tfm.Format(olayput), // feed["publishedAt"],
				"siteID":      feed["siteID"],
				"sitetitle":   feed["sitetitle"],
				"titles":      feed["titles"],
				"url":         feed["url"],
			}
			feedArray = append(feedArray, feedmap)
		}
		return c.JSON(http.StatusOK, map[string][]map[string]interface{}{"data": feedArray})
	}
}

func GetRanking() echo.HandlerFunc {
	return func(c echo.Context) error {
		feed := feedRecord{}
		var feedArray []map[string]interface{}
		// sql02_01 := "SELECT /* sql02_01 */ id, title, URL, image, updateDate, click, siteID FROM articleTBL WHERE id = $1"
		sql02_01 := "SELECT /* sql02_01 */ A.id, A.title, A.URL, A.image, A.updateDate, A.click, S.title FROM articleTBL A INNER JOIN siteTBL S ON A.siteID = S.ID WHERE A.id = $1"

		idsRanking := redis.GetIdsRankingTmp("")
		db := openDB()
		defer db.Close()
		for _, id_count := range idsRanking {
			_id := id_count["id"]
			// if reflect.TypeOf(_id).Kind() != reflect.Int {
			// 	fmt.Println(reflect.TypeOf(_id))
			// 	continue
			// }
			selectFeed := db.QueryRow(sql02_01, _id) //strconv.Itoa(id_count["id"]))
			err := selectFeed.Scan(
				&feed.ID,
				&feed.title,
				&feed.URL,
				&feed.image,
				&feed.updateDate,
				&feed.click,
				// &feed.siteID,
				&feed.siteTitle,
			)
			if err == sql.ErrNoRows {
				continue
			}
			checkError(err)
			feedmap := map[string]interface{}{
				"id":          feed.ID,
				"viewcount":   id_count["viewcount"], //strconv.Itoa(id_count["viewcount"]),
				"titles":      feed.title,
				"url":         feed.URL,
				"image":       feed.image,
				"publishedAt": feed.updateDate,
				"sitetitle":   feed.siteTitle,
			}
			feedArray = append(feedArray, feedmap)
		}
		return c.JSON(http.StatusOK, map[string][]map[string]interface{}{"data": feedArray})
	}
}
