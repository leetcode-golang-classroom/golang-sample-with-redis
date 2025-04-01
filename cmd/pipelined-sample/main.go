package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/leetcode-golang-classroom/golang-sample-with-redis/internal/config"
	"github.com/leetcode-golang-classroom/golang-sample-with-redis/internal/logger"
	myredis "github.com/leetcode-golang-classroom/golang-sample-with-redis/internal/redis"
	"github.com/leetcode-golang-classroom/golang-sample-with-redis/internal/util"
	"github.com/redis/go-redis/v9"
)

func main() {
	jsonLogger := slog.New(slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{
			AddSource: true,
		},
	))
	ctx := logger.CtxWithLogger(context.Background(), jsonLogger)
	config.Init(ctx)
	redisURL := config.AppCfg.RedisUrl

	rdb, err := myredis.New(redisURL)
	if err != nil {
		util.FailOnError(ctx, err, fmt.Sprintf("failed to connect to %s\n", redisURL))
	}
	defer rdb.Close()
	_, err = rdb.Ping(ctx)
	if err != nil {
		util.FailOnError(ctx, err, fmt.Sprintf("failed to ping to %s\n", redisURL))
	}
	// try 10 times
	for i := 0; i < 10; i++ {
		err = rdb.Watch(ctx, func(tx *redis.Tx) error {
			pipe := tx.Pipeline()
			err = pipe.IncrBy(ctx, "p1", 100).Err()
			if err != nil {
				return err
			}
			err = pipe.DecrBy(ctx, "p0", 100).Err()
			if err != nil {
				return err
			}
			_, err = pipe.Exec(ctx)
			return err
		}, "p0")

		if err == nil {
			jsonLogger.Info("transaction execution success")
			break
		} else if err == redis.TxFailedErr {
			jsonLogger.Info("transaction execution failed", slog.Int("i", i))
			continue
		} else {
			util.FailOnError(ctx, err, "failed")
		}
	}

}
