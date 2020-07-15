package dp_redis

import (
	"testing"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	// ServiceName is the name of the service
	ServiceName = "redis"
)

func TestCache_Checker(t *testing.T) {

	Convey("Given that health endpoint returns 'Success'", t, func() {

		// TODO - need to mock this?
		// RedisClient with success health check
		c := &Redis{
			client: NewClient(Options{}).client,
			ttl:    0,
		}

		// CheckState for test validation
		checkState := healthcheck.NewCheckState(ServiceName)

		Convey("Checker updates the CheckState to an OK status", func() {
			c.Checker(checkState)
			So(checkState.Status(), ShouldEqual, healthcheck.StatusOK)
			So(checkState.Message(), ShouldEqual, HealthyMessage)
			So(checkState.StatusCode(), ShouldEqual, 0)
		})
	})

}
