package cache

import "github.com/redis/go-redis/v9"

type RedisClient struct {
	rdb *redis.Client
}

func NewRedisClient(conn string) (*RedisClient, error) {
	opt, err := redis.ParseURL(conn)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	defer client.Close()
	return &RedisClient{
		rdb: client,
	}, nil
}
