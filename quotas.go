package cfclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type QuotaResponse struct {
	Count     int             `json:"total_results"`
	Pages     int             `json:"total_pages"`
	Resources []QuotaResource `json:"resources"`
}

type QuotaResource struct {
	Meta   Meta  `json:"metadata"`
	Entity Quota `json:"entity"`
}

type Quota struct {
	Guid                    string `json:"guid"`
	Name                    string `json:"name"`
	MemoryLimit             int    `json:"memory_limit"`
	InstanceMemoryLimit     int    `json:"instance_memory_limit"`
	TotalRoutes             int    `json:"total_routes"`
	TotalServices           int    `json:"total_services"`
	NonBasicServicesAllowed bool   `json:"non_basic_services_allowed"`
	c                       *Client
}

func (c *Client) GetOrgQuotas() ([]Quota, error) {
	var quotas []Quota
	var quotaResp QuotaResponse
	r := c.NewRequest("GET", "/v2/quota_definitions")
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, fmt.Errorf("Error requesting organization quotas %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading organization quotas request %v", resBody)
	}

	err = json.Unmarshal(resBody, &quotaResp)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling organization quotas %v", err)
	}
	for _, quota := range quotaResp.Resources {
		quota.Entity.Guid = quota.Meta.Guid
		quota.Entity.c = c
		quotas = append(quotas, quota.Entity)
	}
	return quotas, nil
}

func (c *Client) SetOrgQuota(guid string, quotaGuid string) error {
	path := fmt.Sprintf("/v2/organizations/%s", guid)
	r := c.NewRequest("PUT", path)
	r.obj = map[string]interface{}{
		"quota_definition_guid": quotaGuid,
	}
	resp, err := c.DoRequest(r)
	if err != nil {
		return fmt.Errorf("Error setting organization quota %v", err)
	}
	fmt.Printf("%#v\n", resp.Status)

	return nil
}

func (c *Client) GetSpaceQuotas(orgGuid string) ([]Quota, error) {
	var quotas []Quota
	var quotaResp QuotaResponse
	path := fmt.Sprintf("/v2/organizations/%s/space_quota_definitions", orgGuid)
	r := c.NewRequest("GET", path)
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, fmt.Errorf("Error requesting space quotas %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading space quotas request %v", resBody)
	}

	err = json.Unmarshal(resBody, &quotaResp)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling organization quotas %v", err)
	}
	for _, quota := range quotaResp.Resources {
		quota.Entity.Guid = quota.Meta.Guid
		quota.Entity.c = c
		quotas = append(quotas, quota.Entity)
	}
	return quotas, nil
}

func (c *Client) SetSpaceQuota(guid string, quotaGuid string) error {
	path := fmt.Sprintf("/v2/space_quota_definitions/%s/spaces/%s", quotaGuid, guid)
	r := c.NewRequest("PUT", path)
	resp, err := c.DoRequest(r)
	if err != nil {
		return fmt.Errorf("Error setting space quota %v", err)
	}
	fmt.Printf("%#v\n", resp.Status)

	return nil
}

func (c *Client) CreateSpaceQuota(quota Quota, orgGuid string) error {
	r := c.NewRequest("POST", "/v2/space_quota_definitions")
	r.obj = map[string]interface{}{
		"name": quota.Name,
		"non_basic_services_allowed": quota.NonBasicServicesAllowed,
		"total_services":             quota.TotalServices,
		"total_routes":               quota.TotalRoutes,
		"memory_limit":               quota.MemoryLimit,
		"instance_memory_limit":      quota.InstanceMemoryLimit,
		"organization_guid":          orgGuid,
	}
	resp, err := c.DoRequest(r)
	if err != nil {
		return fmt.Errorf("Error creating space quota: %v", err)
	}
	fmt.Printf("%#v\n", resp.Status)

	return nil
}
