package redis

import (
	"fmt"
	_"strconv"
)

func GetIdsRankingTmp(kind string) []map[string]interface{}{
	db := OpenKVS()
	defer db.Close()

	var zsetKey string
	zsetKey_m, zsetKey_w, zsetKey_d := getZsetKeys()
	switch kind {
		case "monthly":
			zsetKey = zsetKey_m
		case "weekly":
			zsetKey = zsetKey_w
		case "daily":
			zsetKey = zsetKey_d
		default: 
			zsetKey = zsetKey_w
	}
	idsranking, err := db.ZRevRangeWithScores(zsetKey, 0, 14).Result()
	checkError(err)
	
	var rankArray []map[string]interface{}
	for _, z := range idsranking {
		/*
		Member_String, isString := z.Member.(string)
		if !isString {
			continue
		}
		_, err := strconv.Atoi(Member_String)
		if err != nil {
			continue
		}
		*/
		memStr, isStr := z.Member.(string)
		if !isStr {
			continue
		}
		if memStr == "null" {
			continue
		}
		rankmap := map[string]interface{} {
			"id":          z.Member,
			"viewcount":   z.Score,
		}
		rankArray = append(rankArray, rankmap)
	}
	fmt.Println(rankArray)
	return rankArray
}