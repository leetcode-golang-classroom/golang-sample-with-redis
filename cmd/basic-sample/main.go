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
	err = rdb.Set(ctx, "key1", "value1", 0)
	if err != nil {
		util.FailOnError(ctx, err, "failed to set key1")
	}
	res, err := rdb.Get(ctx, "key1")
	if err != nil {
		jsonLogger.Error("failed to get key1")
	} else {
		jsonLogger.Info("result get", slog.Any("res", res))
	}
	_, err = rdb.Get(ctx, "key2")
	if err != nil {
		if err == redis.Nil {
			jsonLogger.Error("key2 does not exists")
		} else {
			jsonLogger.Error("err", slog.Any("err", err))
		}
	}
}
