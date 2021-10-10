package elastic

import (
	"fmt"
	"strings"
	"io/ioutil"
)

// https://github.com/elastic/elasticsearch/blob/master/docs/reference/query-dsl/match-query.asciidoc#L268
// https://github.com/elastic/go-elasticsearch/blob/5d2cd13f7d59c109a65bc6cec8c9a238d6c91050/.doc/examples/src/query-dsl-match-query_7f56755fb6c42f7e6203339a6d0cb6e6_test.go
func MatchQueryTest() {
	client := OpenES()

	res, err := client.Search(
		client.Search.WithBody(strings.NewReader(`{
		  "query": {
		    "match": {
		      "titles": "コロナ 自殺"
		    }
		  }
		}`)),
		client.Search.WithPretty(),
	)
	defer res.Body.Close()
	fmt.Println(res, err)
	checkError(err)
}

// https://github.com/elastic/go-elasticsearch/blob/5d2cd13f7d59c109a65bc6cec8c9a238d6c91050/.doc/examples/src/docs-bulk_ae9ccfaa146731ab9176df90670db1c2_test.go
func BulkTest() {
	client := OpenES()

	jsString, err := ioutil.ReadFile("insert_data.js")
	checkError(err)
	res, err := client.Bulk(
		strings.NewReader(string(jsString)),
	)
	defer res.Body.Close()

	fmt.Println(res)
	checkError(err)
}