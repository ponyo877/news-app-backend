package handler

import (
	"context"
	_ "encoding/json"
	"fmt"
	"net/http"
	_ "strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo"
	"github.com/mmcdole/gofeed"
	_ "github.com/ponyo877/news-app-backend/handler/elastic"
	"github.com/ponyo877/news-app-backend/handler/imagectl"
	"github.com/ponyo877/news-app-backend/handler/mongo"
	"github.com/ponyo877/news-app-backend/handler/siever"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
)

func PutPostMongo() echo.HandlerFunc {
	return func(c echo.Context) error {
		update_count := PutPostMongoTmp()
		return c.JSON(http.StatusOK, map[string]int{"update_count": update_count})
	}
}

func PutPostMongoJob() {
	update_count := PutPostMongoTmp()
	fmt.Println("update_count: ", update_count)
}

func PutPostMongoTmp() int {
	ctx := context.Background()
	client := mongo.OpenMongo()
	err := client.Connect(ctx)
	defer client.Disconnect(ctx)
	checkError(err)

	siteinfolist := getSiteInfoListMongo()
	feedparser := gofeed.NewParser()
	feedArray := []interface{}{} // []map[string]interface{}{}
	isDuplicate := map[string]bool{}
	// isVisit := map[int]bool{}
	// update_count := 0
	col := client.Database("newsdb").Collection("article_col")

	for _, siteinfo := range siteinfolist {
		siteID := int(siteinfo["siteID"].(float64))
		siteTitle := siteinfo["sitetitle"].(string)
		latestDate := siteinfo["latestDate"].(string)
		feed, _ := feedparser.ParseURL(siteinfo["rssURL"].(string))
		items := feed.Items
		latesteUpdateDate := items[0].Published

		if latesteUpdateDate > latestDate {
			updateLatestDateMongo(siteID, latesteUpdateDate)
		} else {
			continue
		}
		for _, item := range items {
			if isDuplicate[item.Title] || siever.ContainsNGWord(item.Title) {
				continue
			}
			isDuplicate[item.Title] = true

			if item.Published > latestDate {
				// update_count++
				// siteID := int(siteinfo["siteID"].(float64))
				title_filter := bson.M{"titles": bson.M{"$eq": item.Title}}
				site_filter := bson.M{"siteID": bson.M{"$eq": siteID}}
				filter := bson.M{
					"$and": []interface{}{
						title_filter,
						site_filter,
					},
				}
				_, err = col.DeleteMany(ctx, filter)
				checkError(err)

				feedmap := map[string]interface{}{
					"image":       getImageFromFeedMongo(item.Content),
					"publishedAt": item.Published,
					"sitetitle":   siteTitle,
					"siteID":      siteID,
					"titles":      item.Title,
					"url":         item.Link,
					"acquired":    false,
					"elastic":     false,
				}
				feedArray = append(feedArray, feedmap)
				// if !isVisit[siteID] {
				// 	updateLatestDateMongo(siteID, item.Published)
				// 	isVisit[siteID] = true
				// }
			}
		}
	}
	if len(feedArray) > 0 {
		registerLatestArticleToMongo(feedArray)
	}
	return len(feedArray)
}

func getImageFromFeedMongo(feed string) string {
	reader := strings.NewReader(feed)
	doc, _ := goquery.NewDocumentFromReader(reader)
	imageUrl, _ := doc.Find("img").Attr("src")
	return imagectl.ArrangeImageUrl(imageUrl)
}

// https://qiita.com/h6591/items/f3a7c1bca31cfa634cca
// https://medium.com/since-i-want-to-start-blog-that-looks-like-men-do/%E5%88%9D%E5%BF%83%E8%80%85%E3%81%AB%E9%80%81%E3%82%8A%E3%81%9F%E3%81%84interface%E3%81%AE%E4%BD%BF%E3%81%84%E6%96%B9-golang-48eba361c3b4
// https://noknow.info/it/go/how_to_conveert_between_map_string_interface_and_struct?lang=ja
func registerLatestArticleToMongo(articleList []interface{} /* []map[string]interface{} */) {
	ctx := context.Background()
	client := mongo.OpenMongo()
	err := client.Connect(ctx)
	defer client.Disconnect(ctx)
	checkError(err)
	col := client.Database("newsdb").Collection("article_col")
	_, err = col.InsertMany(ctx, articleList)
	checkError(err)
}

func updateLatestDateMongo(siteID int, updateDate string) {
	ctx := context.Background()
	client := mongo.OpenMongo()
	err := client.Connect(ctx)
	defer client.Disconnect(ctx)
	checkError(err)

	col := client.Database("newsdb").Collection("site_col")
	filter := bson.M{"siteID": bson.M{"$eq": siteID}}
	update := bson.M{"$set": bson.M{"latestDate": updateDate}}
	_, err = col.UpdateOne(ctx, filter, update)
	checkError(err)
}

func getSiteInfoListMongo() []map[string]interface{} {
	ctx := context.Background()
	client := mongo.OpenMongo()
	err := client.Connect(ctx)
	defer client.Disconnect(ctx)
	checkError(err)

	var siteinfolist []map[string]interface{}

	col := client.Database("newsdb").Collection("site_col")
	filter := bson.D{}
	cur, err := col.Find(ctx, filter)
	checkError(err)

	for cur.Next(ctx) {
		var siteinfo bson.M
		err = cur.Decode(&siteinfo)
		fmt.Println(siteinfo)
		checkError(err)
		siteinfolist = append(siteinfolist, siteinfo)
	}
	return siteinfolist
}
