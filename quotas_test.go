package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetOrgQuotas(t *testing.T) {
	Convey("Get organization quotas", t, func() {
		setup(MockRoute{"GET", "/v2/quota_definitions", orgQuotasDefinitionsPayload})
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		quotas, err := client.GetOrgQuotas()
		So(err, ShouldBeNil)

		So(len(quotas), ShouldEqual, 1)
		So(quotas[0].Guid, ShouldEqual, "095a6b8c-31a7-4bc0-a11c-c6a829cfd74c")
		So(quotas[0].Name, ShouldEqual, "default")
	})
}

func TestSetOrgQuota(t *testing.T) {
	Convey("Set organization quota", t, func() {
		setup(MockRoute{"PUT", "/v2/organizations/0f345184-4fd5-4728-8cda-d5ad7f183e9f", orgSetQuotaPayload})
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.SetOrgQuota("0f345184-4fd5-4728-8cda-d5ad7f183e9f", "55b3bbcb-e075-4fc3-904d-733feb8964dc")
		So(err, ShouldBeNil)
	})
}

func TestGetSpaceQuotas(t *testing.T) {
	Convey("Get space quotas", t, func() {
		setup(MockRoute{"GET", "/v2/organizations/9d3a1be7-d504-42bc-987d-90298dcb6b69/space_quota_definitions", spaceQuotasDefinitionsPayload})
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		quotas, err := client.GetSpaceQuotas("9d3a1be7-d504-42bc-987d-90298dcb6b69")
		So(err, ShouldBeNil)

		So(len(quotas), ShouldEqual, 1)
		So(quotas[0].Guid, ShouldEqual, "2203db56-bf87-4ffe-91b9-d46c810fb1b5")
		So(quotas[0].Name, ShouldEqual, "name-1759")
	})
}

func TestSetSpaceQuota(t *testing.T) {
	Convey("Set space quota", t, func() {
		setup(MockRoute{"PUT", "/v2/space_quota_definitions/b8d91aeb-1967-495b-a287-8814ce7bbed0/spaces/5aa4eac1-b034-4261-90ed-d2e6c3ecf76f", spaceSetQuotaPayload})
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.SetSpaceQuota("0f345184-4fd5-4728-8cda-d5ad7f183e9f", "55b3bbcb-e075-4fc3-904d-733feb8964dc")
		So(err, ShouldBeNil)
	})
}

func TestCreateSpaceQuota(t *testing.T) {
	Convey("Create space quota", t, func() {
		setup(MockRoute{"POST", "/v2/space_quota_definitions", spaceCreateQuotaPayload})
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		quota := Quota{
			Guid:                    "17f055b8-b4c8-47cf-8737-0220d5706b4",
			Name:                    "gold_quota",
			MemoryLimit:             5120,
			InstanceMemoryLimit:     -1,
			TotalRoutes:             10,
			TotalServices:           -1,
			NonBasicServicesAllowed: true,
			c: client,
		}

		err = client.CreateSpaceQuota(quota, "c9b4ac17-ab4b-4368-b3e2-5cbf09b17a24")
		So(err, ShouldBeNil)
	})
}
