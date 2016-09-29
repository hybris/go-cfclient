package cfclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type UserResponse struct {
	Count     int            `json:"total_results"`
	Pages     int            `json:"total_pages"`
	Resources []UserResource `json:"resources"`
}

type UserResource struct {
	Meta   Meta `json:"metadata"`
	Entity User `json:"entity"`
}

type User struct {
	Guid            string `json:"guid"`
	Username        string `json:"username"`
	SpaceUrl        string `json:"spaces_url"`
	OrganizationUrl string `json:"organizations_url"`
	Admin           bool   `json:"admin"`
	Active          bool   `json:"active"`
	CreatedAt       string `json:"created_at"`
	c               *Client
}

type UserSummaryResponse struct {
	Guid     string `json:"guid"`
	Username string `json:"username"`
	Active   string `json:"active"`
}

func (c *Client) ListUsers() ([]User, error) {
	var users []User
	var userResp UserResponse
	r := c.NewRequest("GET", "/v2/users")
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, fmt.Errorf("Error requesting users %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading user request %v", resBody)
	}

	err = json.Unmarshal(resBody, &userResp)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling user %v", err)
	}
	for _, user := range userResp.Resources {
		user.Entity.Guid = user.Meta.Guid
		user.Entity.c = c
		users = append(users, user.Entity)
	}
	return users, nil
}
