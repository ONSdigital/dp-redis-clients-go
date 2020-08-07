package redis

import (
	"errors"
	"github.com/ONSdigital/dp-redis/mock"
	"github.com/ONSdigital/dp-sessions-api/session"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClient_Set(t *testing.T) {
	Convey("Given redis client is set up", t, func() {

		mockClient := mock.ClientManagerMock{SetFunc: func(s *session.Session) error {
			return nil
		}}

		Convey("Where the session is valid", func() {
			s := session.NewSession()
			err := mockClient.Set(s)

			Convey("Then the session should be Set in redis", func() {
				So(mockClient.SetCalls(), ShouldHaveLength, 1)
				So(err, ShouldBeNil)
			})
		})

	})

	Convey("Given redis client is set up", t, func() {

		mockClient := mock.ClientManagerMock{SetFunc: func(s *session.Session) error {
			return errors.New("empty session provided")
		}}

		Convey("Where the session is empty", func() {
			var s *session.Session = nil
			err := mockClient.Set(s)

			Convey("Then session is not stored in redis", func() {
				So(mockClient.SetCalls(), ShouldHaveLength, 1)
				So(err, ShouldNotBeEmpty)
			})
		})
	})
}
