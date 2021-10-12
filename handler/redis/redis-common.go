package redis

import (
	"fmt"
	"time"
    "os"
    "strconv"
    "io/ioutil"
    "encoding/json"
    "github.com/go-redis/redis"
)

type KVSConfig struct {
    Host    string  `json:"host"`
    Port    int     `json:"port"`
    Db  	int		`json:"db"`
    Pass    string  `json:"pass"`
}

func checkError(err error) {
	if err != nil {
	        fmt.Fprintf(os.Stderr, "fatal: error: %s", err.Error())
		os.Exit(1)
	}
}


func getZsetKeys() (zsetKey_m string, zsetKey_w string, zsetKey_d string) {
	now := time.Now()

	diff := time.Sunday - now.Weekday()
	Sun := now.AddDate(0, 0, int(diff))

	zsetKey_m = fmt.Sprintf("m_%s-01", now.Format("2006-01"))
	zsetKey_w = fmt.Sprintf("w_%s", Sun.Format("2006-01-02"))
	zsetKey_d = fmt.Sprintf("d_%s", now.Format("2006-01-02"))

	return
}

func OpenKVS() *redis.Client{
	jsonString, err := ioutil.ReadFile("config/config_redis.json")
    checkError(err)
    
    var c KVSConfig
    err = json.Unmarshal(jsonString, &c)
	checkError(err)
	
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Host + ":" + strconv.Itoa(c.Port),
		Password: c.Pass,
		DB:       c.Db,
	})
	return rdb
}

