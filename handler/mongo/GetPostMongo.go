package mongo

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type string_slice_t []string

func (ss string_slice_t) toInt() []int {
	f := make([]int, len(ss))
	for i, v := range ss {
		f[i], _ = strconv.Atoi(v)
	}
	return f
}

// https://qiita.com/h6591/items/a1898bddb6819b27d88f
// https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial
// https://yaruki-strong-zero.hatenablog.jp/entry/go_echo_multiple_param_receive
func GetPostMongoSkip() echo.HandlerFunc {
	return func(c echo.Context) error {
		qfromTmp := c.QueryParam("from")
		qskipIDsTmp := c.QueryParam("skipIDs")

		qfrom, _ := strconv.ParseInt(qfromTmp, 10, 64)
		qskipIDs := string_slice_t(strings.Split(qskipIDsTmp, ",")).toInt()
		fmt.Println("qfrom: ", qfrom)
		fmt.Println("qskipIDs: ", qskipIDs)

		ctx := context.Background()
		client := OpenMongo()
		err := client.Connect(ctx)
		defer client.Disconnect(ctx)
		checkError(err)

		col := client.Database("newsdb").Collection("article_col")

		findOptions := options.Find()
		findOptions.SetSort(bson.D{{"publishedAt", -1}}).SetSkip(qfrom).SetLimit(15)
		var filter interface{}
		if qskipIDsTmp == "" {
			filter = bson.D{}
		} else {
			filter = bson.M{
				"siteID": bson.M{"$nin": qskipIDs},
			}
		}
		cur, err := col.Find(ctx, filter, findOptions)
		checkError(err)

		var feedArray []map[string]interface{}
		for cur.Next(ctx) {
			var feed bson.M
			err = cur.Decode(&feed)
			checkError(err)
			feedArray = append(feedArray, feed)
		}
		return c.JSON(http.StatusOK, map[string][]map[string]interface{}{"data": feedArray})
	}
}

func GetPostMongoLatest() echo.HandlerFunc {
	return func(c echo.Context) error {
		numlimit := int64(150)
		qfirstpublished := c.QueryParam("firstpublished")
		qskipIDsTmp := c.QueryParam("skipIDs")
		qskipIDs := string_slice_t(strings.Split(qskipIDsTmp, ",")).toInt()

		fmt.Println("qfirstpublished: ", qfirstpublished)
		fmt.Println("qskipIDs: ", qskipIDs)

		ctx := context.Background()
		client := OpenMongo()
		err := client.Connect(ctx)
		defer client.Disconnect(ctx)
		checkError(err)

		col := client.Database("newsdb").Collection("article_col")

		findOptions := options.Find()
		findOptions.SetSort(bson.D{{"publishedAt", -1}}).SetLimit(numlimit)

		filter := bson.M{}
		fskipIDsTmp := bson.M{}
		ffirstpublished := bson.M{}
		if qskipIDsTmp != "" {
			fskipIDsTmp = bson.M{
				"siteID": bson.M{"$nin": qskipIDs},
			}
		}
		if qfirstpublished != "" {
			ffirstpublished = bson.M{
				"publishedAt": bson.M{"$gt": qfirstpublished},
			}
		}
		filter = bson.M{
			"$and": []interface{}{
				fskipIDsTmp,
				ffirstpublished,
			},
		}
		cur, err := col.Find(ctx, filter, findOptions)
		checkError(err)

		var feedArray []map[string]interface{}
		var publishedAtTmp string
		ilayout := "2006-01-02T15:04:05-07:00"
		olayput := "2006-01-02 15:04"
		for cur.Next(ctx) {
			var feed bson.M
			err = cur.Decode(&feed)
			// change time format
			publishedAtTmp = feed["publishedAt"].(string)
			tfm, _ := time.Parse(ilayout, publishedAtTmp)
			feed["publishedAt"] = tfm.Format(olayput)

			checkError(err)
			feedArray = append(feedArray, feed)
		}
		// lastpublished := feedArray[numlimit-1]["publishedAt"]
		return c.JSON(http.StatusOK, map[string]interface{}{"data": feedArray})
	}
}

func GetPostMongo() echo.HandlerFunc {
	return func(c echo.Context) error {
		numlimit := int64(15)
		qlastpublished := c.QueryParam("lastpublished")
		qskipIDsTmp := c.QueryParam("skipIDs")
		qskipIDs := string_slice_t(strings.Split(qskipIDsTmp, ",")).toInt()

		fmt.Println("qlastpublished: ", qlastpublished)
		fmt.Println("qskipIDs: ", qskipIDs)

		ctx := context.Background()
		client := OpenMongo()
		err := client.Connect(ctx)
		defer client.Disconnect(ctx)
		checkError(err)

		col := client.Database("newsdb").Collection("article_col")

		findOptions := options.Find()
		findOptions.SetSort(bson.D{{"publishedAt", -1}}).SetLimit(numlimit)

		filter := bson.M{}
		fskipIDsTmp := bson.M{}
		flastpublished := bson.M{}
		if qskipIDsTmp != "" {
			fskipIDsTmp = bson.M{
				"siteID": bson.M{"$nin": qskipIDs},
			}
		}
		if qlastpublished != "" {
			flastpublished = bson.M{
				"publishedAt": bson.M{"$lt": qlastpublished},
			}
		}
		filter = bson.M{
			"$and": []interface{}{
				fskipIDsTmp,
				flastpublished,
			},
		}
		cur, err := col.Find(ctx, filter, findOptions)
		checkError(err)

		var feedArray []map[string]interface{}
		var publishedAtTmp string
		ilayout := "2006-01-02T15:04:05-07:00"
		olayput := "2006-01-02 15:04"
		for cur.Next(ctx) {
			var feed bson.M
			err = cur.Decode(&feed)
			// change time format
			publishedAtTmp = feed["publishedAt"].(string)
			tfm, _ := time.Parse(ilayout, publishedAtTmp)
			feed["publishedAt"] = tfm.Format(olayput)

			checkError(err)
			feedArray = append(feedArray, feed)
		}
		// lastpublished := feedArray[numlimit-1]["publishedAt"]
		return c.JSON(http.StatusOK, map[string]interface{}{"lastpublished": publishedAtTmp, "data": feedArray})
	}
}
