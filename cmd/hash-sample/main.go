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
	var rh1 = RedisHash{
		Name:   "eddie",
		ID:     123,
		Online: true,
	}
	rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, "rh1", "name", rh1.Name)
		pipe.HSet(ctx, "rh1", "id", rh1.ID)
		pipe.HSet(ctx, "rh1", "online", rh1.Online)
		return nil
	})
	var rh2 RedisHash
	err = rdb.HGetAll(ctx, "rh1").Scan(&rh2)
	if err != nil {
		util.FailOnError(ctx, err, "failed on scan rh2")
	}
	jsonLogger.Info("hash sample", slog.Any("rh2", rh2))
}

// RedisHash struct for handle
type RedisHash struct {
	Name   string `redis:"name"`
	ID     int32  `redis:"id"`
	Online bool   `redis:"online"`
}
