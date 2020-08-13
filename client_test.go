package dpredis

import (
	"errors"
	"github.com/ONSdigital/dp-redis/interfaces"
	"testing"
	"time"

	"github.com/ONSdigital/dp-redis/interfaces/mock"
	"github.com/ONSdigital/dp-sessions-api/session"
	"github.com/go-redis/redis"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewClient(t *testing.T) {
	Convey("Given correct redis config", t, func() {
		c, err := setUpClient("123.0.0.1", "1234", 0, 0)

		Convey("Then the client will be created", func() {
			So(err, ShouldBeNil)
			So(c, ShouldNotBeEmpty)
		})
	})

	Convey("Given redis config with missing address", t, func() {
		c, err := setUpClient("", "1234", 0, 0)

		Convey("Then the client will fail to be created", func() {
			So(c, ShouldBeNil)
			So(err, ShouldNotBeEmpty)
			So(err.Error(), ShouldEqual, "address is missing")
		})
	})

	Convey("Given redis config with missing password", t, func() {
		c, err := setUpClient("123.0.0.1", "", 0, 0)

		Convey("Then the client will fail to be created", func() {
			So(c, ShouldBeNil)
			So(err, ShouldNotBeEmpty)
			So(err.Error(), ShouldEqual, "password is missing")
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
				So(err.Error(), ShouldEqual, "failed to store session")
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
				So(err.Error(), ShouldEqual, "session is empty")
			})
		})
	})
}

func setUpClient(addr, password string, database int, ttl time.Duration) (*Client, error){
	c, err := NewClient(Config{
		Addr:     addr,
		Password: password,
		Database: database,
		TTL:      ttl,
	})
	return c, err
}

func setUpMocks(statusCmd redis.StatusCmd) (*mock.RedisClienterMock, *Client) {
	mockRedisClient := &mock.RedisClienterMock{
		PingFunc: nil,
		SetFunc:  func(key string, value interface{}, ttl time.Duration) interfaces.Resulter {
			return &statusCmd
		}}
	return mockRedisClient, &Client{
		client: RedisClient{
			client: mockRedisClient,
		},
		ttl:    0,
	}
}