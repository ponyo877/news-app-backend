package redis

import (
	_ "context"
	"time"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

// var ctx = context.Background()

func IncrViewCount() echo.HandlerFunc {
	return func(c echo.Context) error {
		postIDStr := c.Param("post_id")
		score_m, score_w, score_d := IncrViewCountTmp(postIDStr)
		// return c.JSON(http.StatusOK, map[string]int{"newscore": newscore})
		return c.JSON(http.StatusOK, map[string]int{"score_m": score_m, "score_w": score_w, "score_d": score_d})
	}
}

func GetViewCount() echo.HandlerFunc {
	return func(c echo.Context) error {
		postIDStr := c.Param("post_id")
		score_m, score_w, score_d := GetViewCountTmp(postIDStr)
		// return c.JSON(http.StatusOK, map[string]int{"score": score})
		return c.JSON(http.StatusOK, map[string]int{"score_m": score_m, "score_w": score_w, "score_d": score_d})
	}
}

func IncrViewCountTmp(postID string) (int, int, int) {
	db := OpenKVS()
	defer db.Close()
	// ctx := context.Background()
	// zsetKey := "view_counter"
	zsetKey_m, zsetKey_w, zsetKey_d := getZsetKeys()

	// for Debug
	fmt.Println("👇👇👇 Debug for /redis/put/:post_id 👇👇👇")
	fmt.Printf("zsetKey_m: %v, zsetKey_w: %v, zsetKey_d: %v\n", zsetKey_m, zsetKey_w, zsetKey_d)
	fmt.Printf("Today: %v\n", time.Now().Format("2006-01-02"))

	// newscore, err := db.ZIncrBy(ctx, zsetKey, 1, postID).Result()
	score_m, err := db.ZIncrBy(zsetKey_m, 1, postID).Result()
	score_w, err := db.ZIncrBy(zsetKey_w, 1, postID).Result()
	score_d, err := db.ZIncrBy(zsetKey_d, 1, postID).Result()
	checkError(err)
	return int(score_m), int(score_w), int(score_d)
}

func GetViewCountTmp(postID string) (int, int, int) {
	db := OpenKVS()
	defer db.Close()
	// ctx := context.Background()
	// zsetKey := "view_counter"
	zsetKey_m, zsetKey_w, zsetKey_d := getZsetKeys()

	// score, err := db.ZScore(ctx, zsetKey, postID).Result()
	score_m, err := db.ZScore(zsetKey_m, postID).Result()
	score_w, err := db.ZScore(zsetKey_w, postID).Result()
	score_d, err := db.ZScore(zsetKey_d, postID).Result()
	checkError(err)
	return int(score_m), int(score_w), int(score_d)
}
