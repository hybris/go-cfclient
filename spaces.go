package cfclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type SpaceResponse struct {
	Count     int             `json:"total_results"`
	Pages     int             `json:"total_pages"`
	NextUrl   string          `json:"next_url"`
	Resources []SpaceResource `json:"resources"`
}

type SpaceResource struct {
	Meta   Meta  `json:"metadata"`
	Entity Space `json:"entity"`
}

type Space struct {
	Guid      string        `json:"guid"`
	Name      string        `json:"name"`
	QuotaURL  string        `json:"space_quota_definition_url"`
	QuotaData QuotaResource `json:"quota"`
	OrgURL    string        `json:"organization_url"`
	OrgData   OrgResource   `json:"organization"`
	c         *Client
}

type SpaceSummary struct {
	Guid         string `json:"guid"`
	Name         string `json:"name"`
	ServiceCount int    `json:"service_count"`
	AppCount     int    `json:"app_count"`
	MemDevTotal  int    `json:"mem_dev_total"`
	MemProdTotal int    `json:"mem_prod_total"`
}

type SpaceUsersResponse struct {
	Count     int                  `json:"total_results"`
	Pages     int                  `json:"total_pages"`
	Resources []SpaceUsersResource `json:"resources"`
}

type SpaceUsersResource struct {
	Meta   Meta `json:"metadata"`
	Entity User `json:"entity"`
}

func (s *Space) Org() (Org, error) {
	var orgResource OrgResource
	r := s.c.NewRequest("GET", s.OrgURL)
	resp, err := s.c.DoRequest(r)
	if err != nil {
		return Org{}, fmt.Errorf("Error requesting org %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Org{}, fmt.Errorf("Error reading org request %v", err)
	}

	err = json.Unmarshal(resBody, &orgResource)
	if err != nil {
		return Org{}, fmt.Errorf("Error unmarshaling org %v", err)
	}
	orgResource.Entity.Guid = orgResource.Meta.Guid
	orgResource.Entity.c = s.c
	return orgResource.Entity, nil
}

func (s *Space) Quota() (Quota, error) {
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

func (c *Client) ListSpaces() ([]Space, error) {
	var spaces []Space
	requestUrl := "/v2/spaces"
	for {
		spaceResp, err := c.getSpaceResponse(requestUrl)
		if err != nil {
			return []Space{}, err
		}
		for _, space := range spaceResp.Resources {
			space.Entity.Guid = space.Meta.Guid
			space.Entity.c = c
			spaces = append(spaces, space.Entity)
		}
		requestUrl = spaceResp.NextUrl
		if requestUrl == "" {
			break
		}
	}
	return spaces, nil
}

func (c *Client) getSpaceResponse(requestUrl string) (SpaceResponse, error) {
	var spaceResp SpaceResponse
	r := c.NewRequest("GET", requestUrl)
	resp, err := c.DoRequest(r)
	if err != nil {
		return SpaceResponse{}, fmt.Errorf("Error requesting spaces %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return SpaceResponse{}, fmt.Errorf("Error reading space request %v", err)
	}
	err = json.Unmarshal(resBody, &spaceResp)
	if err != nil {
		return SpaceResponse{}, fmt.Errorf("Error unmarshalling space %v", err)
	}
	return spaceResp, nil
}

func (c *Client) ListSpacesUsers(spaceGuid string, spaceRole string) ([]User, error) {
	var spaceUsers []User
	var spaceUsersResp SpaceUsersResponse
	path := fmt.Sprintf("/v2/spaces/%s/%s", spaceGuid, spaceRole)
	r := c.NewRequest("GET", path)
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, fmt.Errorf("Error requesting space %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading space request %v", resBody)
	}
	err = json.Unmarshal(resBody, &spaceUsersResp)
	if err != nil {
		return nil, fmt.Errorf("Error space organization %v", err)
	}
	for _, user := range spaceUsersResp.Resources {
		user.Entity.Guid = user.Meta.Guid
		spaceUsers = append(spaceUsers, user.Entity)
	}

	return spaceUsers, nil
}

func (c *Client) CreateSpace(spaceName string, orgGuid string, quotaGuid string) error {
	r := c.NewRequest("POST", "/v2/spaces")
	r.obj = map[string]interface{}{
		"name":                  spaceName,
		"organization_guid":     orgGuid,
		"quota_definition_guid": quotaGuid,
	}
	resp, err := c.DoRequest(r)
	if err != nil {
		return fmt.Errorf("Error creating space: %v", err)
	}
	fmt.Printf("%#v\n", resp.Status)

	return nil
}

func (c *Client) SetSpaceUserRole(spaceGuid string, username string, userType string) error {
	path := fmt.Sprintf("/v2/spaces/%s/%s", spaceGuid, userType)
	r := c.NewRequest("PUT", path)
	r.obj = map[string]interface{}{
		"username": username,
	}
	resp, err := c.DoRequest(r)
	if err != nil {
		return fmt.Errorf("Error setting user space role: %v", err)
	}
	fmt.Printf("%#v\n", resp.Status)

	return nil
}
