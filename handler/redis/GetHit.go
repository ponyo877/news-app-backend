package redis

import (
	"fmt"
	_"context"
    "github.com/go-redis/redis"
)

// var ctx = context.Background()

func AddScore() {
	db := OpenKVS()
	defer db.Close()
	// ctx := context.Background()
	zsetKey := "language_rank"
	languages := []*redis.Z{
		&redis.Z{Score: 90.0, Member: "Golang"},
		&redis.Z{Score: 98.0, Member: "Java"},
		&redis.Z{Score: 95.0, Member: "Python"},
		&redis.Z{Score: 97.0, Member: "JavaScript"},
		&redis.Z{Score: 99.0, Member: "C/C++"},
	}
	// languages.append(&redis.Z{Score: 1.0, Member: "Ruby"})

	// _, err := db.ZAdd(ctx, zsetKey, languages...).Result()
	_, err := db.ZAdd(zsetKey, languages...).Result()
	checkError(err)
}

// https://github.com/huruizhi/go_learning_new/blob/master/day11/05redis/main/main.go
func IncrScore(){
	db := OpenKVS()
	defer db.Close()
	zsetKey := "language_rank"

	// newScore, err := db.ZIncrBy(ctx, zsetKey, 10.0, "Ruby").Result()
	newScore, err := db.ZIncrBy(zsetKey, 10.0, "Ruby").Result()
	checkError(err)
	fmt.Printf("Golang's score is %f now.\n", newScore)
}