package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	
	diff := time.Sunday - now.Weekday()
	this_monday := now.AddDate(0, 0, int(diff))

	redis_monthly := fmt.Sprintf("%s-01", now.Format("2006-01"))
	redis_weekly := this_monday.Format("2006-01-02")
	redis_daily := now.Format("2006-01-02")

	fmt.Printf("redis_monthly: %s\n", redis_monthly)
	fmt.Printf("redis_weekly: %s\n", redis_weekly)
	fmt.Printf("redis_daily: %s\n", redis_daily)
}
