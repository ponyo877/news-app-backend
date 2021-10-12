package mongo

import (
	"fmt"
    "os"
    "strconv"
    "io/ioutil"
    "encoding/json"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
    Host    string  `json:"host"`
    Port    int     `json:"port"`
}

func checkError(err error) {
	if err != nil {
	        fmt.Fprintf(os.Stderr, "fatal: error: %s", err.Error())
		os.Exit(1)
	}
}

func OpenMongo() *mongo.Client{
	jsonString, err := ioutil.ReadFile("config/config_mongo.json")
    checkError(err)
    
    var c MongoConfig
    err = json.Unmarshal(jsonString, &c)
	checkError(err)
	
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + c.Host + ":" + strconv.Itoa(c.Port)))
	checkError(err)
	return client
}