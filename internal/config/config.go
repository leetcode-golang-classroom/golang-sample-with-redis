package config

import (
	"context"

	"github.com/leetcode-golang-classroom/golang-sample-with-redis/internal/logger"
	"github.com/leetcode-golang-classroom/golang-sample-with-redis/internal/util"
	"github.com/spf13/viper"
)

type Config struct {
	RedisUrl string `mapstructure:"REDIS_URL"`
}

var AppCfg *Config

func Init(ctx context.Context) {
	log := logger.FromContext(ctx)
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigType("env")
	v.SetConfigName(".env")
	v.AutomaticEnv()
	// bind environment
	util.FailOnError(ctx, v.BindEnv("REDIS_URL"), "failed to bind env REDIS_URL")
	err := v.ReadInConfig()
	if err != nil {
		log.Warn("Load from environment variable")
	}
	err = v.Unmarshal(&AppCfg)
	if err != nil {
		util.FailOnError(ctx, err, "Failed to read enivronment")
	}
}
