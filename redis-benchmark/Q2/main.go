package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v9"
)

var counts = []int{
	10000,
	50000,
	100000,
	200000,
	500000,
}

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Username: "",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Panic(err)
	}

	for _, n := range counts {
		log.Printf("[count]: %v\n", n)

		// flush db before insert
		if err := rdb.FlushDB(ctx).Err(); err != nil {
			log.Panic(err)
		}

		// log memory info before
		before, err := rdb.Info(ctx, "memory").Result()
		if err != nil {
			log.Panic(err)
		}
		log.Printf("[before] info memory: %v\n", before)

		// insert data to redis concurrently
		wg := &sync.WaitGroup{}
		wg.Add(n)
		for i := 0; i < n; i++ {
			go func() {
				defer wg.Done()
				key := fmt.Sprint("key_%d_%d", n, i)
				if err := rdb.Set(ctx, key, i, 0).Err(); err != nil {
					log.Printf("set error: %v", err)
				}
			}()
		}
		wg.Wait()

		// log memory info after
		after, err := rdb.Info(ctx, "memory").Result()
		if err != nil {
			log.Fatal(err.Error())
		}

		// log result in file
		log.Printf("[after] info memory: %v\n\n", after)

		// flush db after insert
		rdb.FlushDB(ctx)
		time.Sleep(5 * time.Second)
	}
}
