package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListSpaces(t *testing.T) {
	Convey("List Space", t, func() {
		mocks := []MockRoute{
			{"GET", "/v2/spaces", listSpacesPayload},
			{"GET", "/v2/spacesPage2", listSpacesPayloadPage2},
		}
		setupMultiple(mocks)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		spaces, err := client.ListSpaces()
		So(err, ShouldBeNil)

		So(len(spaces), ShouldEqual, 4)
		So(spaces[0].Guid, ShouldEqual, "8efd7c5c-d83c-4786-b399-b7bd548839e1")
		So(spaces[0].Name, ShouldEqual, "dev")
		So(spaces[1].Guid, ShouldEqual, "657b5923-7de0-486a-9928-b4d78ee24931")
		So(spaces[1].Name, ShouldEqual, "demo")
		So(spaces[2].Guid, ShouldEqual, "9ffd7c5c-d83c-4786-b399-b7bd54883977")
		So(spaces[2].Name, ShouldEqual, "test")
		So(spaces[3].Guid, ShouldEqual, "329b5923-7de0-486a-9928-b4d78ee24982")
		So(spaces[3].Name, ShouldEqual, "prod")
	})
}

func TestSpaceOrg(t *testing.T) {
	Convey("Find space org", t, func() {
		mocks := []MockRoute{
			{"GET", "/v2/space_quota_definitions/19f507dd-69d4-4f76-a4fb-ca6b7cff4956", spaceQuotaPayload},
			{"GET", "/v2/org/foobar", orgPayload},
		}
		setupMultiple(mocks)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		space := &Space{
			Guid:     "123",
			Name:     "test space",
			OrgURL:   "/v2/org/foobar",
			QuotaURL: "/v2/space_quota_definitions/19f507dd-69d4-4f76-a4fb-ca6b7cff4956",
			c:        client,
		}
		org, err := space.Org()
		So(err, ShouldBeNil)

		So(org.Name, ShouldEqual, "test-org")
		So(org.Guid, ShouldEqual, "da0dba14-6064-4f7a-b15a-ff9e677e49b2")

		quota, err := space.Quota()
		So(err, ShouldBeNil)

		So(quota.Name, ShouldEqual, "name-1491")
		So(quota.Guid, ShouldEqual, "a9097bc8-c6cf-4a8f-bc47-623fa22e8019")
	})
}

func TestListSpaceUsers(t *testing.T) {
	Convey("Find space users", t, func() {
		setup(MockRoute{"GET", "/v2/spaces/foobar/developers", spaceUsersPayload})
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		spaceUsers, err := client.ListSpacesUsers("foobar", "developers")
		So(err, ShouldBeNil)

		So(len(spaceUsers), ShouldEqual, 1)
		So(spaceUsers[0].Guid, ShouldEqual, "138557be-69b8-4c71-8ba1-502bb2d0dd29")
		So(spaceUsers[0].Username, ShouldEqual, "user1")
	})
}

func TestCreateSpace(t *testing.T) {
	Convey("Create Space", t, func() {
		setup(MockRoute{"POST", "/v2/spaces", createSpacePayload})
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.CreateSpace("foobar", "8efd7c5c-d83c-4786-b399-b7bd548839e1", "657b5923-7de0-486a-9928-b4d78ee24931")
		So(err, ShouldBeNil)
	})
}

func TestSetSpaceUserRole(t *testing.T) {
	Convey("Set Space User role", t, func() {
		setup(MockRoute{"PUT", "/v2/spaces/8efd7c5c-d83c-4786-b399-b7bd548839e1/developers", setSpaceUserRolePayload})
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.SetSpaceUserRole("8efd7c5c-d83c-4786-b399-b7bd548839e1", "name-2142", "developers")
		So(err, ShouldBeNil)
	})
}
