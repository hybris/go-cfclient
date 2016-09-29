package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListOrgs(t *testing.T) {
	Convey("List Org", t, func() {
		setup(MockRoute{"GET", "/v2/organizations", listOrgsPayload})
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		orgs, err := client.ListOrgs()
		So(err, ShouldBeNil)

		So(len(orgs), ShouldEqual, 2)
		So(orgs[0].Guid, ShouldEqual, "a537761f-9d93-4b30-af17-3d73dbca181b")
		So(orgs[0].Name, ShouldEqual, "demo")
	})
}

func TestOrgSpaces(t *testing.T) {
	Convey("Get spaces by org", t, func() {
		setup(MockRoute{"GET", "/v2/organizations/foo/spaces", orgSpacesPayload})
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		spaces, err := client.OrgSpaces("foo")
		So(err, ShouldBeNil)

		So(len(spaces), ShouldEqual, 1)
		So(spaces[0].Guid, ShouldEqual, "b8aff561-175d-45e8-b1e7-67e2aedb03b6")
		So(spaces[0].Name, ShouldEqual, "test")
	})
}

func TestOrgSummary(t *testing.T) {
	Convey("Get organization summary", t, func() {
		setup(MockRoute{"GET", "/v2/organizations/foo/summary", orgSummaryPayload})
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		summary, err := client.OrgSummary("foo")
		So(err, ShouldBeNil)

		So(len(summary), ShouldEqual, 1)
		So(summary[0].Guid, ShouldEqual, "cf63a51f-9dcd-43f8-8552-a08a4a3b4df3")
		So(summary[0].Name, ShouldEqual, "name-1375")
		So(summary[0].ServiceCount, ShouldEqual, 0)
		So(summary[0].AppCount, ShouldEqual, 0)
		So(summary[0].MemDevTotal, ShouldEqual, 0)
		So(summary[0].MemProdTotal, ShouldEqual, 0)
	})
}

func TestOrgQuota(t *testing.T) {
	Convey("Get organization quota", t, func() {
		setup(MockRoute{"GET", "/v2/quota_definitions/80f3e539-a8c0-4c43-9c72-649df53da8cb", orgQuotaPayload})
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		org := &Org{
			Guid:     "0c69f181-2d31-4326-ac33-be2b114a5f99",
			Name:     "foobar",
			QuotaURL: "/v2/quota_definitions/80f3e539-a8c0-4c43-9c72-649df53da8cb",
			c:        client,
		}
		quota, err := org.Quota()
		So(err, ShouldBeNil)

		So(quota.Name, ShouldEqual, "name-1996")
		So(quota.Guid, ShouldEqual, "80f3e539-a8c0-4c43-9c72-649df53da8cb")
		So(quota.TotalServices, ShouldEqual, 60)
		So(quota.TotalRoutes, ShouldEqual, 1000)
		So(quota.MemoryLimit, ShouldEqual, 20480)
		So(quota.NonBasicServicesAllowed, ShouldEqual, true)
		So(quota.InstanceMemoryLimit, ShouldEqual, 1024)
	})
}

func TestOrgCreate(t *testing.T) {
	Convey("Create organization", t, func() {
		setup(MockRoute{"POST", "/v2/organizations", orgCreatePayload})
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.CreateOrg("my-org-name", "b7887f5c-34bb-40c5-9778-577572e4fb2d")
		So(err, ShouldBeNil)
	})
}

func TestSetOrgUserRole(t *testing.T) {
	Convey("Set organization user role", t, func() {
		mocks := []MockRoute{
			{"PUT", "/v2/organizations/272b3566-04bd-4f3a-b83d-75deb8a67649/auditors", orgUserRolePayload},
			{"PUT", "/v2/organizations/272b3566-04bd-4f3a-b83d-75deb8a67649/users", orgUserPayload},
		}
		setupMultiple(mocks)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.SetOrgUserRole("272b3566-04bd-4f3a-b83d-75deb8a67649", "name-1766", "auditors")
		So(err, ShouldBeNil)
	})
}
