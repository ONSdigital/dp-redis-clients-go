package redis

//go:generate moq -out mock/mockresulter.go -pkg mock . Resulter
//go:generate moq -out mock/mockclienter.go -pkg mock . Clienter

import (
	"context"
	"errors"
	"time"

	"github.com/ONSdigital/dp-sessions-api/session"
	"github.com/ONSdigital/log.go/log"
	goredis "github.com/go-redis/redis"
)

var ctx = context.Background()

//Resulter - interface for redis.StatsCMD
type Resulter interface {
	Result() (string, error)
}

//Clienter - interface for redis
type Clienter interface {
	Set(string, string, time.Duration) Resulter
	Ping() Resulter
}

//GoRedisClient - structure for the redis client
type GoRedisClient struct {
	client *goredis.Client
}

func (rc *GoRedisClient) Set(key string, value string, ttl time.Duration) Resulter {
	return rc.client.Set(key, value, ttl)
}

func (rc GoRedisClient) Ping() Resulter {
	return rc.client.Ping()
}

// Client - structure for the cache client
type Client struct {
	client GoRedisClient
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
		client:  GoRedisClient{client:goredis.NewClient(&goredis.Options{
			Addr: c.Addr,
			Password: c.Password,
			DB: c.Database,
		})},
		ttl:     0,
	}, nil
}

//Set add session to redis
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
	}
	log.Event(ctx, msg, log.INFO)

	return nil
}
