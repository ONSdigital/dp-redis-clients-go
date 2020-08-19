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

	if c.TTL == 0 {
		return nil, errors.New("zero is not a valid ttl")
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
		log.Event(ctx, "session is empty", log.ERROR)
		return errors.New("session is empty")
	}

	sJSON, err := s.MarshalJSON()
	if err != nil {
		log.Event(ctx, "failed to marshal session", log.Error(err), log.ERROR)
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
func (c *Client) GetByID(id string) (*session.Session, error) {
	if id == "" {
		log.Event(ctx, "id value is blank", log.ERROR)
		return nil, errors.New("id value is blank")
	}

	msg, err := c.client.Get(id).Result()
	if err != nil {
		log.Event(ctx, msg, log.Error(err), log.ERROR)
		return nil, err
	}

	var s *session.Session

	err = json.Unmarshal([]byte(msg), &s)
	if err != nil {
		log.Event(ctx, "failed to unmarshal session", log.Error(err), log.ERROR)
		return nil, err
	}

	return s, nil
}

// Set - redis implementation of Set
func (rc *RedisClient) Set(key string, value string, expiration time.Duration) *redis.StatusCmd {
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