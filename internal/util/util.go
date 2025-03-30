package util

import (
	"context"
	"fmt"
	"os"

	"github.com/leetcode-golang-classroom/golang-sample-with-redis/internal/logger"
)

func FailOnError(ctx context.Context, err error, msg string) {
	log := logger.FromContext(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("%s: %s", msg, err))
		os.Exit(1)
	}
}
