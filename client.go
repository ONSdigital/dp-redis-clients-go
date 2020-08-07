package redis

import (
	"context"
	"errors"
	"github.com/ONSdigital/log.go/log"
	"time"

	"github.com/ONSdigital/dp-sessions-api/session"
	goredis "github.com/go-redis/redis"
)

var ctx = context.Background()

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

	msg, err := c.client.Set(s.ID, json, c.ttl).Result()
	if err != nil {
		log.Event(ctx, msg, log.Error(err), log.ERROR)
	}
	log.Event(ctx, msg, log.INFO)

	return nil
}
