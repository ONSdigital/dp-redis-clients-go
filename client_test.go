package dpredis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ONSdigital/dp-redis/interfaces/mock"
	"github.com/ONSdigital/dp-sessions-api/session"
	"github.com/go-redis/redis"
	. "github.com/smartystreets/goconvey/convey"
)

const testTTL = 30 * time.Minute

var (
	resp = []byte(`{"id":"1234","email":"user@email.com","start":"2020-08-13T08:40:18.652Z","last_accessed":"2020-08-13T08:40:18.652Z"}`)
	ctx  = context.Background()
)

func TestNewClient(t *testing.T) {
	Convey("Given correct redis config", t, func() {
		c, err := NewClient(Config{
			Addr:     "123.0.0.1",
			Password: "1234",
			Database: 0,
			TTL:      testTTL,
		})

		Convey("Then the client will be created", func() {
			So(err, ShouldBeNil)
			So(c, ShouldNotBeEmpty)
		})
	})

	Convey("Given redis config with missing address", t, func() {
		c, err := NewClient(Config{
			Addr:     "",
			Password: "1234",
			Database: 0,
			TTL:      testTTL,
		})

		Convey("Then the client will fail to be created", func() {
			So(c, ShouldBeNil)
			So(err, ShouldNotBeEmpty)
			So(err.Error(), ShouldEqual, "address is missing")
		})
	})

	Convey("Given redis config with missing password", t, func() {
		c, err := NewClient(Config{
			Addr:     "123.0.0.1",
			Password: "",
			Database: 0,
			TTL:      testTTL,
		})

		Convey("Then the client will fail to be created", func() {
			So(c, ShouldBeNil)
			So(err, ShouldNotBeEmpty)
			So(err.Error(), ShouldEqual, "password is missing")
		})
	})

	Convey("Given redis config with a zero ttl", t, func() {
		c, err := NewClient(Config{
			Addr:     "123.0.0.1",
			Password: "1234",
			Database: 0,
			TTL:      0,
		})

		Convey("Then the client will fail to be created", func() {
			So(c, ShouldBeNil)
			So(err, ShouldNotBeEmpty)
			So(err.Error(), ShouldEqual, "zero is not a valid ttl")
		})
	})
}

func TestClient_Set(t *testing.T) {
	Convey("Given a valid session", t, func() {
		mockRedisClient, client := setUpMocks(*redis.NewStatusResult("success", nil), *redis.NewStringCmd())

		s := &session.Session{
			ID:           "1234",
			Email:        "user@email.com",
			Start:        time.Now(),
			LastAccessed: time.Now(),
		}

		Convey("When cache attempts to store session", func() {
			err := client.Set(ctx, s)

			Convey("Then session should be stored", func() {
				So(err, ShouldBeNil)
				So(mockRedisClient.SetCalls(), ShouldHaveLength, 1)
			})
		})
	})

	Convey("Given a valid session", t, func() {
		mockRedisClient, client := setUpMocks(*redis.NewStatusResult("fail", errors.New("failed to store session")), *redis.NewStringCmd())

		s := &session.Session{
			ID:           "1234",
			Email:        "user@email.com",
			Start:        time.Now(),
			LastAccessed: time.Now(),
		}

		Convey("When cache attempts to store session", func() {
			err := client.Set(ctx, s)

			Convey("Then session will not be stored", func() {
				So(mockRedisClient.SetCalls(), ShouldHaveLength, 1)
				So(err, ShouldNotBeEmpty)
				So(err.Error(), ShouldEqual, "failed to store session")
			})
		})
	})

	Convey("Given an invalid session", t, func() {
		mockRedisClient, client := setUpMocks(*redis.NewStatusCmd(), *redis.NewStringCmd())

		var s *session.Session = nil

		Convey("When cache attempts to store session", func() {
			err := client.Set(ctx, s)

			Convey("Then session will not be stored", func() {
				So(mockRedisClient.SetCalls(), ShouldHaveLength, 0)
				So(err, ShouldNotBeEmpty)
				So(err.Error(), ShouldEqual, ErrEmptySession)
			})
		})
	})
}

func TestClient_GetByID(t *testing.T) {
	Convey("Given a session ID", t, func() {
		mockRedisClient, client := setUpMocks(*redis.NewStatusCmd(), *redis.NewStringResult(string(resp), nil))

		Convey("When cache attempts to get session", func() {
			s, err := client.GetByID(ctx, "1234")
			So(err, ShouldBeNil)

			Convey("Then session is returned", func() {
				So(mockRedisClient.GetCalls(), ShouldHaveLength, 1)
				So(s, ShouldNotBeEmpty)
				So(s.ID, ShouldEqual, "1234")
			})
		})
	})

	Convey("Given a blank session ID", t, func() {
		mockRedisClient, client := setUpMocks(*redis.NewStatusCmd(), *redis.NewStringCmd())

		Convey("When cache attempts to get session", func() {
			s, err := client.GetByID(ctx, "")

			Convey("Then session is not returned", func() {
				So(mockRedisClient.GetCalls(), ShouldHaveLength, 0)
				So(s, ShouldBeNil)
				So(err, ShouldNotBeEmpty)
				So(err.Error(), ShouldEqual, ErrEmptySessionID)
			})
		})
	})

	Convey("Given a session ID", t, func() {
		mockRedisClient, client := setUpMocks(*redis.NewStatusCmd(), *redis.NewStringResult("", errors.New("unexpected end of JSON input")))

		Convey("When cache attempts to get session", func() {
			s, err := client.GetByID(ctx,"1234")

			Convey("Then session is not returned", func() {
				So(mockRedisClient.GetCalls(), ShouldHaveLength, 1)
				So(s, ShouldBeNil)
				So(err, ShouldNotBeEmpty)
				So(err.Error(), ShouldEqual, "unexpected end of JSON input")
			})
		})
	})
}

func setUpMocks(statusCmd redis.StatusCmd, stringCmd redis.StringCmd) (*mock.RedisClienterMock, *Client) {
	mockRedisClient := &mock.RedisClienterMock{
		PingFunc: nil,
		SetFunc: func(key string, value interface{}, ttl time.Duration) *redis.StatusCmd {
			return &statusCmd
		},
		GetFunc: func(key string) *redis.StringCmd {
			return &stringCmd
		}}
	return mockRedisClient, &Client{
		client: RedisClient{
			client: mockRedisClient,
		},
		ttl: testTTL,
	}
}
