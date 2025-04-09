package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var ctx = context.Background()

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(addr string) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       1,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}

	return &RedisClient{Client: client}
}

func (r *RedisClient) StoreRefreshToken(userID string, token string, expiry time.Duration) error {
	return r.Client.Set(ctx, "refresh:"+userID, token, expiry).Err()
}

func (r *RedisClient) GetRefreshToken(userID string) (string, error) {
	return r.Client.Get(ctx, "refresh:"+userID).Result()
}

func (r *RedisClient) DeleteKey(key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r *RedisClient) BlacklistAccessToken(token string, expiry time.Duration) error {
	return r.Client.Set(ctx, "blacklist:"+token, "true", expiry).Err()
}

func (r *RedisClient) IsTokenBlacklisted(token string) (bool, error) {
	exists, err := r.Client.Exists(ctx, "blacklist:"+token).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}
