package dpredis

import (
	"errors"
	"testing"
	"time"

	"github.com/ONSdigital/dp-redis/interfaces/mock"
	"github.com/ONSdigital/dp-sessions-api/session"
	"github.com/go-redis/redis"
	. "github.com/smartystreets/goconvey/convey"
)

const testTTL = 30 * time.Minute

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
			So(err, ShouldEqual, ErrEmptyAddress)
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
			So(err, ShouldEqual, ErrEmptyPassword)
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
			So(err, ShouldEqual, ErrInvalidTTL)
		})
	})
}

func TestClient_Set(t *testing.T) {
	Convey("Given a valid session", t, func() {
		mockRedisClient, mockClient := setUpMocks(*redis.NewStatusResult("success", nil))

		s := &session.Session{
			ID:           "1234",
			Email:        "user@email.com",
			Start:        time.Now(),
			LastAccessed: time.Now(),
		}

		Convey("When cache attempts to store session", func() {
			err := mockClient.Set(s)

			Convey("Then session should be stored", func() {
				So(err, ShouldBeNil)
				So(mockRedisClient.SetCalls(), ShouldHaveLength, 1)
			})
		})
	})

	Convey("Given a valid session", t, func() {
		mockRedisClient, mockClient := setUpMocks(*redis.NewStatusResult("fail", errors.New("failed to store session")))

		s := &session.Session{
			ID: "1234",
			Email: "user@email.com",
			Start: time.Now(),
			LastAccessed: time.Now(),
		}

		Convey("When cache attempts to store session", func() {
			err := mockClient.Set(s)

			Convey("Then session will not be stored", func() {
				So(mockRedisClient.SetCalls(), ShouldHaveLength, 1)
				So(err, ShouldNotBeEmpty)
				So(err.Error(), ShouldEqual, "redis client.Set returned an unexpected error: failed to store session")
			})
		})
	})

	Convey("Given an  invalid session", t, func() {
		mockRedisClient, mockClient := setUpMocks(*redis.NewStatusCmd())

		var s *session.Session = nil

		Convey("When cache attempts to store session", func() {
			err := mockClient.Set(s)

			Convey("Then session will not be stored", func() {
				So(mockRedisClient.SetCalls(), ShouldHaveLength, 0)
				So(err, ShouldNotBeEmpty)
				So(err, ShouldEqual, ErrEmptySession)
			})
		})
	})
}

func setUpMocks(statusCmd redis.StatusCmd) (*mock.RedisClienterMock, *Client) {
	mockRedisClient := &mock.RedisClienterMock{
		PingFunc: nil,
		SetFunc:  func(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
			return &statusCmd
		}}
	return mockRedisClient, &Client{
		client: RedisClient{
			client: mockRedisClient,
		},
		ttl:    0,
	}
}