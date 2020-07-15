package dp_redis

import (
	"github.com/go-redis/redis"
	"time"
)

// Redis - structure for the redis client
type Redis struct {
	client *redis.Client
	ttl    time.Duration
}

// Options - config options for the redis client
type Options struct {
	Addr     string
	Password string
	DB       int
	ttl      time.Duration
}

// NewClient - returns new redis client with provided config options
func NewClient(o Options) *Redis {
	return &Redis{
		client: redis.NewClient(&redis.Options{
			Addr: o.Addr,
			Password: o.Password,
			DB: o.DB,
		}),
		ttl:    o.ttl,
	}
}
