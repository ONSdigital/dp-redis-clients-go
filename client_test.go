package dpredis

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/ONSdigital/dp-redis/interfaces/mock"
	"github.com/ONSdigital/dp-sessions-api/session"
	"github.com/go-redis/redis"
	. "github.com/smartystreets/goconvey/convey"
)

var resp = []byte(`{"id":"1234","email":"user@email.com","start":"2020-08-13T08:40:18.652Z","last_accessed":"2020-08-13T08:40:18.652Z"}`)

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
		mockRedisClient, mockClient := setUpMocks(*redis.NewStatusResult("success", nil), *redis.NewStringCmd())

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
		mockRedisClient, mockClient := setUpMocks(*redis.NewStatusResult("fail", errors.New("failed to store session")), *redis.NewStringCmd())

		s := &session.Session{
			ID:           "1234",
			Email:        "user@email.com",
			Start:        time.Now(),
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

	Convey("Given an invalid session", t, func() {
		mockRedisClient, mockClient := setUpMocks(*redis.NewStatusCmd(), *redis.NewStringCmd())

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

func TestClient_GetByID(t *testing.T) {
	Convey("Given a session ID", t, func() {
		mockRedisClient, mockClient := setUpMocks(*redis.NewStatusCmd(), *redis.NewStringResult(string(resp), nil))

		Convey("When cache attempts to get session", func() {
			s, err := mockClient.GetByID("1234")
			So(err, ShouldBeNil)

			msg, err := mockRedisClient.Get(s.ID).Result()
			So(err, ShouldBeNil)

			var storedSession session.Session
			err = json.Unmarshal([]byte(msg), &storedSession)
			So(err, ShouldBeNil)

			Convey("Then session is returned", func() {
				So(storedSession.ID, ShouldEqual, s.ID)
				So(storedSession.Email, ShouldEqual, s.Email)
			})
		})
	})

	Convey("Given a blank session ID", t, func() {
		_, mockClient := setUpMocks(*redis.NewStatusCmd(), *redis.NewStringCmd())

		Convey("When cache attempts to get session", func() {
			_, err := mockClient.GetByID("")

			Convey("Then session is returned", func() {
				So(err, ShouldNotBeEmpty)
				So(err.Error(), ShouldEqual, "id value is blank")
			})
		})
	})

	Convey("Given a session ID", t, func() {
		mockRedisClient, _ := setUpMocks(*redis.NewStatusCmd(), *redis.NewStringResult("", errors.New("")))

		Convey("When cache attempts to get session", func() {
			msg, err := mockRedisClient.Get("1234").Result()

			var storedSession session.Session
			err = json.Unmarshal([]byte(msg), &storedSession)

			Convey("Then session is returned", func() {
				So(err, ShouldNotBeEmpty)
				So(err.Error(), ShouldEqual, "unexpected end of JSON input")
			})
		})
	})
}

func setUpClient(addr, password string, database int, ttl time.Duration) (*Client, error) {
	c, err := NewClient(Config{
		Addr:     addr,
		Password: password,
		Database: database,
		TTL:      ttl,
	})
	return c, err
}

func setUpMocks(statusCmd redis.StatusCmd, stringCmd redis.StringCmd) (*mock.RedisClienterMock, *Client) {
	mockRedisClient := &mock.RedisClienterMock{
		PingFunc: nil,
		SetFunc: func(key string, value interface{}, ttl time.Duration) *redis.StatusCmd {
			return &statusCmd
		},
		GetFunc: func(id string) *redis.StringCmd {
			return &stringCmd
		}}
	return mockRedisClient, &Client{
		client: RedisClient{
			client: mockRedisClient,
		},
		ttl: 0,
	}
}
