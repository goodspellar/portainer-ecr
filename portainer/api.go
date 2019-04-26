package portainer

import (
	"encoding/json"
	"log"
	"strconv"
)

type Registry struct {
	ID             int    `json:"Id,omitempty"`
	Name           string `json:"Name,omitempty"`
	URL            string `json:"URL,omitempty"`
	Authentication bool   `json:"Authentication,omitempty"`
	Username       string `json:"Username,omitempty"`
	Password       string `json:"Password,omitempty"`
}

func (c *Client) GetRegistries() []Registry {
	var registries []Registry
	resp, err := c.CallAPI("GET", "/api/registries", nil)
	if err != nil {
		log.Fatal("Could not get registries from Portainer: ", err)
	}
	err = json.NewDecoder(resp.Body).Decode(&registries)
	if err != nil {
		log.Fatal("Failed to parse registry response from Portainer: ", err)
	}

	return registries
}

func (c *Client) UpdateRegistry(r *Registry) {
	id := strconv.Itoa(r.ID)

	resp, err := c.CallAPI("PUT", "/api/registries/"+id, *r)
	if err != nil {
		log.Fatal("Could not update Portainer ECR registry: ", err)
	}

	if resp.StatusCode/100 != 2 {
		log.Fatal("Problem updating Portainer registry")
	}

	log.Println("Successfully updated registry: ", r.Name)
}
