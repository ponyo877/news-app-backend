package mongo

import (
    "context"
    "net/http"
    "github.com/labstack/echo"
    "go.mongodb.org/mongo-driver/bson"
    _"go.mongodb.org/mongo-driver/mongo/options"
)

func GetSiteInfoMongo() echo.HandlerFunc {
    return func(c echo.Context) error {
        ctx := context.Background()
        client := OpenMongo()
        err := client.Connect(ctx)
        defer client.Disconnect(ctx)
		checkError(err)
		
		col := client.Database("newsdb").Collection("site_col")
        filter := bson.D{}
		cur, err := col.Find(ctx, filter)
        checkError(err)
        
        var feedArray []map[string]interface{}
		for cur.Next(ctx) {
            var feed bson.M
			err = cur.Decode(&feed);
			checkError(err)
            feedArray = append(feedArray, feed)
        }
        return c.JSON(http.StatusOK, map[string][]map[string]interface{}{"data": feedArray})
    }
}