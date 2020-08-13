package dpredis

import (
	"context"
	"errors"
	"time"

	. "github.com/ONSdigital/dp-redis/interfaces"
	"github.com/ONSdigital/dp-sessions-api/session"
	"github.com/ONSdigital/log.go/log"
	"github.com/go-redis/redis"
)

var ctx = context.Background()

// RedisClient - structure for the redis client
type RedisClient struct {
	client RedisClienter
}

// Client - structure for the cache client
type Client struct {
	client RedisClient
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
		client:  RedisClient{client:redis.NewClient(&redis.Options{
			Addr: c.Addr,
			Password: c.Password,
			DB: c.Database,
		})},
		ttl:     0,
	}, nil
}

// Set - add session to redis
func (c *Client) Set(s *session.Session) error {
	if s == nil {
		log.Event(ctx, "session is empty", log.ERROR)
		return errors.New("session is empty")
	}

	json, err := s.MarshalJSON()
	if err != nil {
		log.Event(ctx, "failed to marshal session", log.Error(err), log.ERROR)
		return err
	}

	msg, err := c.client.Set(s.ID, string(json), c.ttl).Result()
	if err != nil {
		log.Event(ctx, msg, log.Error(err), log.ERROR)
		return err
	}
	log.Event(ctx, msg, log.INFO)

	return nil
}

// Set - redis implementation of Set
func (rc *RedisClient) Set(key string, value string, ttl time.Duration) *redis.StatusCmd {
	return rc.client.Set(key, value, ttl)
}

// Ping - redis implementation of Ping
func (rc *RedisClient) Ping() Resulter {
	return rc.client.Ping()
}