package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/ponyo877/news-app-backend/handler/elastic"
	"github.com/ponyo877/news-app-backend/handler/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// https://www.elastic.co/jp/blog/the-go-client-for-elasticsearch-working-with-data
// https://kb.objectrocket.com/mongo-db/how-to-update-many-mongodb-documents-using-the-golang-driver-447
func RegisterToSearchEngine() {
	var countSuccessful uint64
	indexName := "searchdb"
	ctx := context.Background()
	mongoclient := mongo.OpenMongo()
	err := mongoclient.Connect(ctx)
	defer mongoclient.Disconnect(ctx)
	checkError(err)

	col := mongoclient.Database("newsdb").Collection("article_col")
	cur, err := col.Find(ctx, bson.M{"elastic": false})
	checkError(err)

	es := elastic.OpenES()
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  indexName, // The default index name
		Client: es,        // The Elasticsearch client
		// NumWorkers:    numWorkers,       // The number of worker goroutines
		// FlushBytes:    int(flushBytes),  // The flush threshold in bytes
		// FlushInterval: 30 * time.Second, // The periodic flush interval
	})

	for cur.Next(ctx) {
		var feed bson.M
		err = cur.Decode(&feed)
		checkError(err)

		feed["id"] = feed["_id"]
		delete(feed, "_id")
		b, err := json.Marshal(feed)
		checkError(err)

		objID := feed["id"].(primitive.ObjectID).Hex()
		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: objID,
				Body:       bytes.NewReader(b),

				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&countSuccessful, 1)
				},
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						fmt.Printf("ERROR: %s", err)
					} else {
						fmt.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)
		checkError(err)
	}
	_, err = col.UpdateMany(ctx, bson.M{"elastic": false}, bson.M{"$set": bson.M{"elastic": true}})
	checkError(err)

	err = bi.Close(context.Background())
	checkError(err)
}
