package elastic

import (
	"strings"
	"encoding/json"
	"time"
	"fmt"
	_"go.mongodb.org/mongo-driver/bson/primitive"
)

// for Code Refactoring
// https://infraya.work/posts/go_json_parse_aws/

func GetSearchResultTmp(searchwords string) []map[string]interface{}{
	client := OpenES()
	// query := "{\"query\": { \"match\": { \"titles\": \"" + searchwords + "\"}}}"
	query := fmt.Sprintf(
`{"query": {"match": {"titles": "%s"}}}`, // , "fields": ["title", "body"]}}}`,
		searchwords)
	res, err := client.Search(
		client.Search.WithBody(strings.NewReader(query)),
		client.Search.WithPretty(),
	)
	checkError(err)
	defer res.Body.Close()
	var jsonstrings map[string]map[string][]map[string]map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&jsonstrings)
	prosjsonstrings := jsonstrings["hits"]["hits"]
	var searchArray []map[string]interface{}

	var publishedAtTmp string
	// var idTmp string
	ilayout := "2006-01-02T15:04:05-07:00"
	olayput := "2006-01-02 15:04"

	for _, searchmaptmp := range prosjsonstrings {
		// change time format
		feed := searchmaptmp["_source"]
		publishedAtTmp = feed["publishedAt"].(string)
		tfm, _ := time.Parse(ilayout, publishedAtTmp)
		feed["publishedAt"] = tfm.Format(olayput)
		feed["_id"] = feed["id"].(string) // (primitive.ObjectID).Hex()
		searchArray = append(searchArray, feed)
	}
	
	return searchArray
}