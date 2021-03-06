package elastic

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/elastic/go-elasticsearch/v7"
)

type ESConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: error: %s", err.Error())
		os.Exit(1)
	}
}

func OpenES() *elasticsearch.Client {
	jsonString, err := ioutil.ReadFile("config/config_es.json")
	checkError(err)

	var c ESConfig
	err = json.Unmarshal(jsonString, &c)
	checkError(err)

	cfg := elasticsearch.Config{
		Addresses: []string{
			// "http://localhost:9200",
			"https://" + c.Host + ":" + strconv.Itoa(c.Port),
			// "http://" + c.Host + ":" + strconv.Itoa(c.Port),
		},
		Username: c.User,
		Password: c.Pass,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	checkError(err)
	return es
}
