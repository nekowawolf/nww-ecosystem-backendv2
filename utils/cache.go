package utils

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nekowawolf/airdropv2/config"
)

func GetOrSetCache[T any](key string, ttl time.Duration, fetchFunc func() (T, error)) (T, error) {
	var result T
	ctx := context.Background()

	cachedData, err := config.RedisClient.Get(ctx, key).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedData), &result); err == nil {
			return result, nil
		}
	}

	result, err = fetchFunc()
	if err != nil {
		return result, err
	}

	if bytes, err := json.Marshal(result); err == nil {
		config.RedisClient.Set(ctx, key, string(bytes), ttl)
	}

	return result, nil
}

func InvalidateCache(keys ...string) {
	if len(keys) == 0 {
		return
	}
	ctx := context.Background()
	config.RedisClient.Del(ctx, keys...)
}

func InvalidateCachePrefix(prefix string) {
	ctx := context.Background()
	iter := config.RedisClient.Scan(ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		config.RedisClient.Del(ctx, iter.Val())
	}
}