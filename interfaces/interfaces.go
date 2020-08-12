package interfaces

//go:generate moq -out mock/mock_redisclienter.go -pkg mock . RedisClienter

import (
	"time"

	"github.com/go-redis/redis"
)

// RedisClienter - interface for redis
type RedisClienter interface {
	Set(string, interface{}, time.Duration) *redis.StatusCmd
	Ping() *redis.StatusCmd
}
