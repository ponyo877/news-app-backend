package elastic

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/elastic/go-elasticsearch"
)

// https://github.com/elastic/go-elasticsearch/blob/5d2cd13f7d59c109a65bc6cec8c9a238d6c91050/.doc/examples/src/query-dsl-match-query_7f56755fb6c42f7e6203339a6d0cb6e6_test.go

var (
	_ = fmt.Printf
	_ = os.Stdout
	_ = elasticsearch.NewDefaultClient
)

// <https://github.com/elastic/elasticsearch/blob/master/docs/reference/query-dsl/match-query.asciidoc#L268>
//
// --------------------------------------------------------------------------------
// GET /_search
// {
//    "query": {
//        "match" : {
//            "message": {
//                "query" : "ny city",
//                "auto_generate_synonyms_phrase_query" : false
//            }
//        }
//    }
// }
// --------------------------------------------------------------------------------

func match_query_test() {
	client := OpenES()

	res, err := client.Search(
		es.Search.WithBody(strings.NewReader(`{
		  "query": {
		    "match": {
		      "message": {
		        "query": "ny city",
		        "auto_generate_synonyms_phrase_query": false
		      }
		    }
		  }
		}`)),
		es.Search.WithPretty(),
	)
	defer res.Body.Close()
	fmt.Println(res, err)
	checkError(err)
}