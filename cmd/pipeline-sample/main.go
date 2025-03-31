package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

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
	err = rdb.Set(ctx, "t1", time.Now().UTC(), 0)
	if err != nil {
		util.FailOnError(ctx, err, "failed to set t1")
	}
	pipe := rdb.Pipeline(ctx)
	t1 := pipe.Get(ctx, "t1")
	jsonLogger.Info(fmt.Sprintf("pipe 執行前的 t1=%v", t1))
	for i := 0; i < 10; i++ {
		pipe.Set(ctx, fmt.Sprintf("p%v", i), i, 0)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		util.FailOnError(ctx, err, "failed to pipeline")
	}
	jsonLogger.Info(fmt.Sprintf("pipe 執行後的 t1=%v", t1))

	cmds, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for i := 0; i < 10; i++ {
			pipe.Get(ctx, fmt.Sprintf("p%v", i))
		}
		return nil
	})
	if err != nil {
		util.FailOnError(ctx, err, "failed to pipeline")
	}
	for i, cmd := range cmds {
		jsonLogger.Info(fmt.Sprintf("p%v=%v", i, cmd.(*redis.StringCmd).Val()))
	}
}
