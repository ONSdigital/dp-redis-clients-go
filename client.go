package redis

import (
	"errors"
	"time"

	goredis "github.com/go-redis/redis"
)

// Client - structure for the redis client
type Client struct {
	client *goredis.Client
	ttl    time.Duration
}

// Config - config options for the redis client
type Config struct {
	Addr     string
	Password string
	Database int
	TTL      time.Duration
}

// NewClient - returns new redis client with provided config options
func NewClient(c Config) (*Client, error) {
	if c.Addr == "" {
		return nil, errors.New("address is missing")
	}

	if c.Password == "" {
		return nil, errors.New("password is missing")
	}

	return &Client{
		client: goredis.NewClient(&goredis.Options{
			Addr:     c.Addr,
			Password: c.Password,
			DB:       c.Database,
		}),
		ttl: c.TTL,
	}, nil
}
