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

var (
	resp = []byte(`{"id":"1234","email":"user@email.com","start":"2020-08-13T08:40:18.652Z","last_accessed":"2020-08-13T08:40:18.652Z"}`)
)

func TestNewClient(t *testing.T) {
	Convey("Given NewClient returns new redis client", t, func() {

		Convey("When correct redis configuration is provided", func() {
			c, err := NewClient(Config{
				Addr:     "123.0.0.1",
				Password: "1234",
				Database: 0,
				TTL:      testTTL,
			})

			Convey("Then a new redis client will be returned with no error", func() {
				So(err, ShouldBeNil)
				So(c, ShouldNotBeEmpty)
			})
		})

	})

	Convey("Given NewClient returns an error", t, func() {

		Convey("When the redis configurations address is empty", func() {
			c, err := NewClient(Config{
				Addr:     "",
				Password: "1234",
				Database: 0,
				TTL:      testTTL,
			})

			Convey("Then the client will not be created and the empty address error is returned", func() {
				So(c, ShouldBeNil)
				So(err, ShouldNotBeEmpty)
				So(err, ShouldEqual, ErrEmptyAddress)
			})
		})

	})

	Convey("Given NewClient returns an error", t, func() {

		Convey("When the redis configurations password is empty", func() {
			c, err := NewClient(Config{
				Addr:     "123.0.0.1",
				Password: "",
				Database: 0,
				TTL:      testTTL,
			})

			Convey("Then the client will not be created and the empty password error is returned", func() {
				So(c, ShouldBeNil)
				So(err, ShouldNotBeEmpty)
				So(err, ShouldEqual, ErrEmptyPassword)
			})
		})

	})

	Convey("Given NewClient returns an error", t, func() {

		Convey("When the redis configurations ttl is zero", func() {
			c, err := NewClient(Config{
				Addr:     "123.0.0.1",
				Password: "1234",
				Database: 0,
				TTL:      0,
			})

			Convey("Then the client will not be created and the invalid ttl error is returned", func() {
				So(c, ShouldBeNil)
				So(err, ShouldNotBeEmpty)
				So(err, ShouldEqual, ErrInvalidTTL)
			})
		})
	})
}

func TestClient_Set(t *testing.T) {
	Convey("Given a valid sessions and redis client.Set returns no error", t, func() {
		mockRedisClient, client := setUpMocks(*redis.NewStatusResult("success", nil), *redis.NewStringCmd())


		Convey("When there is a valid session", func() {
			s := &session.Session{
				ID:           "1234",
				Email:        "user@email.com",
				Start:        time.Now(),
				LastAccessed: time.Now(),
			}

			err := client.Set(s)

			Convey("Then the session is stored in the cache and no error is returned", func() {
				So(err, ShouldBeNil)
				So(mockRedisClient.SetCalls(), ShouldHaveLength, 1)
			})
		})
	})

	Convey("Given a valid session and redis client.Set returns an error", t, func() {
		mockRedisClient, client := setUpMocks(*redis.NewStatusResult("fail", errors.New("failed to store session")), *redis.NewStringCmd())

		Convey("When there is a valid session but redis client.Set errors ", func() {
			s := &session.Session{
				ID:           "1234",
				Email:        "user@email.com",
				Start:        time.Now(),
				LastAccessed: time.Now(),
			}

			err := client.Set(s)

			Convey("Then the session will not be stored in the cache and an error is returned", func() {
				So(mockRedisClient.SetCalls(), ShouldHaveLength, 1)
				So(err, ShouldNotBeEmpty)
				So(err.Error(), ShouldEqual, "failed to store session")
			})
		})
	})

	Convey("Given an invalid session and redis client.Set returns an error", t, func() {
		mockRedisClient, client := setUpMocks(*redis.NewStatusCmd(), *redis.NewStringCmd())

		Convey("When there is an invalid session", func() {
			var s *session.Session = nil
			err := client.Set(s)

			Convey("Then the session will not be stored in the cache and an error is returned", func() {
				So(mockRedisClient.SetCalls(), ShouldHaveLength, 0)
				So(err, ShouldNotBeEmpty)
				So(err, ShouldEqual, ErrEmptySession)
			})
		})
	})
}

func TestClient_GetByID(t *testing.T) {
	Convey("Given a session ID client.GetByID returns a session", t, func() {
		mockRedisClient, client := setUpMocks(*redis.NewStatusCmd(), *redis.NewStringResult(string(resp), nil))

		Convey("When client uses the ID to get the session", func() {
			s, err := client.GetByID("1234")
			So(err, ShouldBeNil)

			Convey("Then redis client.Get is called and returns the session", func() {
				So(mockRedisClient.GetCalls(), ShouldHaveLength, 1)
				So(s, ShouldNotBeEmpty)
				So(s.ID, ShouldEqual, "1234")
			})
		})
	})

	Convey("Given a blank session ID client.GetByID returns an error", t, func() {
		mockRedisClient, client := setUpMocks(*redis.NewStatusCmd(), *redis.NewStringCmd())

		Convey("When client.GetByID is called has an empty ID", func() {
			s, err := client.GetByID("")

			Convey("Then client.GetByID returns an error and no session is returned", func() {
				So(mockRedisClient.GetCalls(), ShouldHaveLength, 0)
				So(s, ShouldBeNil)
				So(err, ShouldNotBeEmpty)
				So(err, ShouldEqual, ErrEmptySessionID)
			})
		})
	})

	Convey("Given a session ID client.GetByID returns an error", t, func() {
		mockRedisClient, client := setUpMocks(*redis.NewStatusCmd(), *redis.NewStringResult("", errors.New("unexpected end of JSON input")))

		Convey("When client.GetByID is called with a valid session ID", func() {
			s, err := client.GetByID("1234")

			Convey("Then the redis client.Get returns an error and no session is returned", func() {
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
