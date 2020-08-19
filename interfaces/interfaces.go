package interfaces

//go:generate moq -out mock/mock_redisclienter.go -pkg mock . RedisClienter

import (
	"time"

	"github.com/go-redis/redis"
)

// RedisClienter - interface for redis
type RedisClienter interface {
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(string) *redis.StringCmd
	Ping() *redis.StatusCmd
}
