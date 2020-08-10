package redis

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
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