package cfclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type OrgResponse struct {
	Count     int           `json:"total_results"`
	Pages     int           `json:"total_pages"`
	Resources []OrgResource `json:"resources"`
}

type OrgResource struct {
	Meta   Meta `json:"metadata"`
	Entity Org  `json:"entity"`
}

type Org struct {
	Guid      string        `json:"guid"`
	Name      string        `json:"name"`
	QuotaURL  string        `json:"quota_definition_url"`
	QuotaData QuotaResource `json:"quota"`
	c         *Client
}

type OrgSummaryResponse struct {
	Guid   string         `json:"guid"`
	Name   string         `json:"name"`
	Status string         `json:"status"`
	Spaces []SpaceSummary `json:"spaces"`
}

func (c *Client) ListOrgs() ([]Org, error) {
	var orgs []Org
	var orgResp OrgResponse
	r := c.NewRequest("GET", "/v2/organizations")
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, fmt.Errorf("Error requesting organizations %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading organization request %v", resBody)
	}

	err = json.Unmarshal(resBody, &orgResp)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling organization %v", err)
	}
	for _, org := range orgResp.Resources {
		org.Entity.Guid = org.Meta.Guid
		org.Entity.c = c
		orgs = append(orgs, org.Entity)
	}
	return orgs, nil
}

func (c *Client) OrgSpaces(guid string) ([]Space, error) {
	var spaces []Space
	var spaceResp SpaceResponse
	path := fmt.Sprintf("/v2/organizations/%s/spaces", guid)
	r := c.NewRequest("GET", path)
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, fmt.Errorf("Error requesting space %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading space request %v", resBody)
	}

	err = json.Unmarshal(resBody, &spaceResp)
	if err != nil {
		return nil, fmt.Errorf("Error space organization %v", err)
	}
	for _, space := range spaceResp.Resources {
		space.Entity.Guid = space.Meta.Guid
		space.Entity.c = c
		spaces = append(spaces, space.Entity)
	}

	return spaces, nil

}

func (c *Client) OrgSummary(guid string) ([]SpaceSummary, error) {
	var orgSummaryResponse OrgSummaryResponse

	path := fmt.Sprintf("/v2/organizations/%s/summary", guid)
	r := c.NewRequest("GET", path)
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, fmt.Errorf("Error requesting space %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading space request %v", resBody)
	}

	err = json.Unmarshal(resBody, &orgSummaryResponse)
	if err != nil {
		return nil, fmt.Errorf("Error space organization summary %v", err)
	}

	return orgSummaryResponse.Spaces, nil

}

func (s *Org) Quota() (Quota, error) {
	var quotaResource QuotaResource
	r := s.c.NewRequest("GET", s.QuotaURL)
	resp, err := s.c.DoRequest(r)
	if err != nil {
		return Quota{}, fmt.Errorf("Error requesting quota %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Quota{}, fmt.Errorf("Error reading quota request %v", err)
	}

	err = json.Unmarshal(resBody, &quotaResource)
	if err != nil {
		return Quota{}, fmt.Errorf("Error unmarshaling quota %v", err)
	}
	quotaResource.Entity.Guid = quotaResource.Meta.Guid
	quotaResource.Entity.c = s.c
	return quotaResource.Entity, nil
}

func (c *Client) CreateOrg(orgName string, quotaGuid string) error {
	r := c.NewRequest("POST", "/v2/organizations")
	r.obj = map[string]interface{}{
		"name":                  orgName,
		"quota_definition_guid": quotaGuid,
	}
	resp, err := c.DoRequest(r)
	if err != nil {
		return fmt.Errorf("Error creating organization: %v", err)
	}
	fmt.Printf("%#v\n", resp.Status)

	return nil
}

func (c *Client) SetOrgUserRole(orgGuid string, username string, userType string) error {
	path := fmt.Sprintf("/v2/organizations/%s/%s", orgGuid, userType)
	r := c.NewRequest("PUT", path)
	r.obj = map[string]interface{}{
		"username": username,
	}
	resp, err := c.DoRequest(r)
	if err != nil {
		return fmt.Errorf("Error setting user org role: %v", err)
	}
	fmt.Printf("%#v\n", resp.Status)

	path = fmt.Sprintf("/v2/organizations/%s/users", orgGuid)
	r = c.NewRequest("PUT", path)
	r.obj = map[string]interface{}{
		"username": username,
	}
	resp, err = c.DoRequest(r)
	if err != nil {
		return fmt.Errorf("Error setting user in org users: %v", err)
	}
	fmt.Printf("%#v\n", resp.Status)

	return nil
}
