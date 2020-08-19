package dpredis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	. "github.com/ONSdigital/dp-redis/interfaces"
	"github.com/ONSdigital/dp-sessions-api/session"
	"github.com/ONSdigital/log.go/log"
	"github.com/go-redis/redis"
)

const (
	ErrEmptySessionID    = "session id required but was empty"
	ErrEmptySession      = "session is empty"
	ErrFailedToUnmarshal = "failed to unmarshal get session json response"
	ErrFailedToMarshal   = "failed to marshal session"
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
		return nil, errors.New("address is missing")
	}

	if c.Password == "" {
		return nil, errors.New("password is missing")
	}

	if c.TTL == 0 {
		return nil, errors.New("zero is not a valid ttl")
	}

	return &Client{
		client: RedisClient{client: redis.NewClient(&redis.Options{
			Addr:     c.Addr,
			Password: c.Password,
			DB:       c.Database,
		})},
		ttl: c.TTL,
	}, nil
}

// Set - add session to redis
func (c *Client) Set(ctx context.Context, s *session.Session) error {
	if s == nil {
		log.Event(ctx, ErrEmptySession, log.ERROR)
		return errors.New(ErrEmptySession)
	}

	sJSON, err := s.MarshalJSON()
	if err != nil {
		log.Event(ctx, ErrFailedToMarshal, log.Error(err), log.ERROR)
		return err
	}

	msg, err := c.client.Set(s.ID, sJSON, c.ttl).Result()
	if err != nil {
		log.Event(ctx, msg, log.Error(err), log.ERROR)
		return err
	}

	return nil
}

// GetByID - gets a session from redis using its ID
func (c *Client) GetByID(ctx context.Context, id string) (*session.Session, error) {
	if id == "" {
		log.Event(ctx, ErrEmptySessionID, log.ERROR)
		return nil, errors.New(ErrEmptySessionID)
	}

	msg, err := c.client.Get(id).Result()
	if err != nil {
		log.Event(ctx, msg, log.Error(err), log.ERROR)
		return nil, err
	}

	var s *session.Session

	err = json.Unmarshal([]byte(msg), &s)
	if err != nil {
		log.Event(ctx, ErrFailedToUnmarshal, log.Error(err), log.ERROR)
		return nil, err
	}

	return s, nil
}

// Set - redis implementation of Set
func (rc *RedisClient) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return rc.client.Set(key, value, expiration)
}

// Get - redis implementation of Get
func (rc *RedisClient) Get(key string) *redis.StringCmd {
	return rc.client.Get(key)
}

// Ping - redis implementation of Ping
func (rc *RedisClient) Ping() *redis.StatusCmd {
	return rc.client.Ping()
}
