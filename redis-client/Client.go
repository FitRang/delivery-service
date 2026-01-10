package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	rdb *redis.Client
}

func NewRedisClient(conn string) (*RedisClient, error) {
	opt, err := redis.ParseURL(conn)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisClient{
		rdb: client,
	}, nil
}

func (r *RedisClient) Close() error {
	return r.rdb.Close()
}
