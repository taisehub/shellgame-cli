package infrastructure

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/taise-hub/shellgame-cli/interfaces"
)

var (
	options redis.Options = redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
)

type redisHandler struct {
	client *redis.Client
}

func NewRedisHandler() interfaces.RedisHandler {
	cli := redis.NewClient(&options)
	h := new(redisHandler)
	h.client = cli
	return h
}

func (h *redisHandler) ListGet(ctx context.Context, key string) ([]string, error) {
	return h.client.LRange(ctx, key, 0, -1).Result()
}

func (h *redisHandler) ListSet(ctx context.Context, key string, value []string) error {
	return h.client.RPush(ctx, key, value).Err()
}
