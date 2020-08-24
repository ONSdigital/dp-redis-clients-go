package dpredis

import (
	"errors"
	"fmt"
	"time"

	. "github.com/ONSdigital/dp-redis/interfaces"
	"github.com/ONSdigital/dp-sessions-api/session"
	"github.com/go-redis/redis"
)

var (
	ErrEmptySession   = errors.New("session required but was empty")
	ErrEmptyAddress   = errors.New("redis host address required but was empty")
	ErrEmptyPassword  = errors.New("redis password required but was empty")
	ErrInvalidTTL     = errors.New("redis client ttl cannot be 0")
)

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
		return nil, ErrEmptyAddress
	}

	if c.Password == "" {
		return nil, ErrEmptyPassword
	}

	if c.TTL == 0 {
		return nil, ErrInvalidTTL
	}

	return &Client{
		client:  RedisClient{client:redis.NewClient(&redis.Options{
			Addr: c.Addr,
			Password: c.Password,
			DB: c.Database,
		})},
		ttl:     c.TTL,
	}, nil
}

// Set - add session to redis
func (c *Client) Set(s *session.Session) error {
	if s == nil {
		return ErrEmptySession
	}

	json, err := s.MarshalJSON()
	if err != nil {
		return err
	}

	err = c.client.Set(s.ID, string(json), c.ttl).Err()
	if err != nil {
		return fmt.Errorf("redis client.Set returned an unexpected error: %w", err)
	}

	return nil
}

// Set - redis implementation of Set
func (rc *RedisClient) Set(key string, value string, expiration time.Duration) *redis.StatusCmd {
	return rc.client.Set(key, value, expiration)
}

// Ping - redis implementation of Ping
func (rc *RedisClient) Ping() *redis.StatusCmd {
	return rc.client.Ping()
}