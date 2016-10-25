package main

import (
	"log"

	"github.com/ekaley/client"
)

// Config contains the ipam API configuration parameters
// required to connect the ipam Go client to the API.
type Config struct {
	Address string
	Scheme  string
}

// Client creates a ipam API client which is
// utilized by the terraform.Provider.
func (c *Config) Client() (*api.Client, error) {
	client := api.NewClient(c.Address)

	log.Printf("Ipam Client Configured: %s\n", c.Address)

	return client, nil
}
