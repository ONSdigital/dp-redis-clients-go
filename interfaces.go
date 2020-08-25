package dpredis

//go:generate moq -out mock_redisclienter.go . RedisClienter

import (
	"time"

	"github.com/go-redis/redis"
)

// RedisClienter - interface for redis
type RedisClienter interface {
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(key string) *redis.StringCmd
	FlushAll() * redis.StatusCmd
	Ping() *redis.StatusCmd
}
