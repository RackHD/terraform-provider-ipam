package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("IPAM_ADDRESS", "127.0.0.1:8000"),
			},

			"scheme": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("IPAM_SCHEME", "http"),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"ipam_pool":        resourceIpamPool(),
			"ipam_subnet":      resourceIpamSubnet(),
			"ipam_reservation": resourceIpamReservation(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Address: d.Get("address").(string),
		Scheme:  d.Get("scheme").(string),
	}

	return config.Client()
}
