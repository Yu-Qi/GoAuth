package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	r "github.com/zeromicro/go-zero/core/stores/redis"

	"github.com/Yu-Qi/GoAuth/pkg/config"
)

// Client .
var Client *redis.Client

// LockClient .
var LockClient *r.Redis

// RedisOptions connection options used by redis
var RedisOptions *redis.Options

// ErrorRedisNil .
var ErrorRedisNil = "redis: nil"

// ErrorRedisZsetEmpty .
var ErrorRedisZsetEmpty = "zset is empty"

func init() {
	redisHost := config.GetString("REDIS_HOST")
	redisPort := config.GetString("REDIS_PORT")

	if redisHost != "" && redisPort != "" {
		password := config.GetString("REDIS_AUTH")

		RedisOptions = &redis.Options{
			Addr:     redisHost + ":" + redisPort,
			Password: password,
		}
		Client = redis.NewClient(RedisOptions)
		LockClient = r.New(redisHost+":"+redisPort, r.WithPass(password))

		return
	}
	panic("missing redis config")
}

// Exists check if key exists
func Exists(ctx context.Context, key string) bool {
	if Client == nil {
		panic("redis client is nil")
	}
	cmd := Client.Exists(ctx, key)
	if cmd == nil {
		return false
	}
	return (*cmd).Val() > 0
}

// Expire set expire time for key
func Expire(ctx context.Context, key string, ttl time.Duration) error {
	if Client == nil {
		panic("redis client is nil")
	}
	cmd := Client.Expire(ctx, key, ttl)
	if cmd == nil {
		return nil
	}
	return (*cmd).Err()
}

// Get get value from redis
func Get(ctx context.Context, key string) (interface{}, error) {
	if Client == nil {
		panic("redis client is nil")
	}
	cmd := Client.Get(ctx, key)
	if cmd == nil {
		return nil, fmt.Errorf("redis: nil")
	}
	return (*cmd).Result()
}

// Set set string to redis
func Set(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	if Client == nil {
		panic("redis client is nil")
	}
	cmd := Client.Set(ctx, key, value, duration)
	if cmd == nil {
		return nil
	}

	return (*cmd).Err()
}

// SetWithObject set object to redis
func SetWithObject(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	objectJson, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if Client == nil {
		panic("redis client is nil")
	}
	cmd := Client.Set(ctx, key, objectJson, duration)
	if cmd == nil {
		return nil
	}

	return (*cmd).Err()
}
