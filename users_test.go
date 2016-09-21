package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListUsers(t *testing.T) {
	Convey("List User", t, func() {
		setup(MockRoute{"GET", "/v2/users", listUsersPayload})
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		users, err := client.ListUsers()
		So(err, ShouldBeNil)

		So(len(users), ShouldEqual, 2)
		So(users[0].Guid, ShouldEqual, "08582a96-cbef-463c-822e-bda8d4284cc7")
		So(users[0].Username, ShouldEqual, "demo")
	})
}
