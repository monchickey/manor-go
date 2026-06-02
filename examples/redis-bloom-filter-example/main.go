package main

import (
	"context"
	"fmt"
	"time"

	"github.com/monchickey/manor-go/v2/bloomfilter"
	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Username: "redis",
		Password: "redis-123456",
		DB:       0,
	})
	defer rdb.Close()

	redisBf := bloomfilter.NewRedisBloomFilter(rdb)

	tsbf, err := bloomfilter.NewTimeSeriesBucketBloomFilter(redisBf, "test:bf", 0.000001, 10000000, "hour", time.Hour)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	now_before := now.Add(-5 * time.Hour)
	now_after := now.Add(5 * time.Hour)
	row1 := bloomfilter.TimeSeriesRow{EventTime: &now, EventId: "user1"}
	row2 := bloomfilter.TimeSeriesRow{EventTime: &now, EventId: "user2"}
	row3 := bloomfilter.TimeSeriesRow{EventTime: &now, EventId: "user3"}
	row4 := bloomfilter.TimeSeriesRow{EventTime: &now_before, EventId: "user4"}
	row5 := bloomfilter.TimeSeriesRow{EventTime: &now_after, EventId: "user5"}

	err = tsbf.AddBatch([]bloomfilter.TimeSeriesRow{row1, row2, row3, row4, row5})
	if err != nil {
		panic(err)
	}

	exists, err := tsbf.ExistsBatch([]bloomfilter.TimeSeriesRow{row1, row2, row3,
		{EventTime: &now, EventId: "user111"}, row4, row5, {EventTime: &now_after, EventId: "user222"}})
	if err != nil {
		panic(err)
	}

	fmt.Println(exists)

	key := tsbf.GetBucket(row1.EventTime)
	expectedKey := "test:bf:" + key
	count, _ := rdb.Exists(context.Background(), expectedKey).Result()
	if count != 1 {
		panic(fmt.Sprintf("expected key %s to exist, got %d", expectedKey, count))
	}

	// 检查 TTL 是否生效
	// ttl, _ := rdb.TTL(context.Background(), expectedKey).Result()
	// if ttl < time.Hour-5*time.Minute || ttl > time.Hour+5*time.Minute {
	// 	panic(fmt.Sprintf("expected TTL near 24h, got %v", ttl))
	// }

	startTime := time.Now()

	// 十万分之一， 1000000
	// 100 批量 10000 次：2.54s、1.92s
	// 1000 批量 100000 次：171.1s、114s
	// 100 批量 100000 次：33.9s、31s

	// 百万分之一，10000000
	// 100 批量 100000 次：39s、40s

	rows := make([]bloomfilter.TimeSeriesRow, 100)
	for i := 0; i < 100; i++ {
		rows[i] = bloomfilter.TimeSeriesRow{
			EventTime: &now,
			EventId:   fmt.Sprintf("user%d", i),
		}
	}

	fmt.Printf("生成数据用时: %f\n", time.Since(startTime).Seconds())

	startTime = time.Now()

	for i := 0; i < 100000; i++ {
		err := tsbf.AddBatch(rows)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("添加用时: %f\n", time.Since(startTime).Seconds())

	startTime = time.Now()

	for range 100000 {
		exists, err = tsbf.ExistsBatch(rows)
		if err != nil {
			panic(err)
		}
		if !allTrue(exists) {
			panic("exists")
		}
	}

	fmt.Printf("判断用时: %f\n", time.Since(startTime).Seconds())

}

func allTrue(arr []bool) bool {
	for _, v := range arr {
		if !v {
			return false
		}
	}
	return true
}
