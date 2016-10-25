package main

import (
	"fmt"
	"log"

	"github.com/RackHD/ipam/resources"
	"github.com/ekaley/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIpamSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpamSubnetsCreate,
		Read:   resourceIpamSubnetsShow,
		Update: resourceIpamSubnetsUpdate,
		Delete: resourceIpamSubnetsDelete,

		Schema: map[string]*schema.Schema{
			"uuid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"metadata": &schema.Schema{
				Type:     schema.TypeMap,
				ForceNew: true,
				Optional: true,
			},
			"pool": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"start": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"end": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceIpamSubnetsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	p := newSubnetFromResource(d)

	subnet, err := client.Subnets().CreateShowSubnet(d.Get("pool").(string), p)
	if err != nil {
		return fmt.Errorf("Failed to create Subnet: %s\n", err)
	}
	d.SetId(subnet.ID)
	d.Set("uuid", subnet.ID)
	log.Println(d.Id())
	log.Printf("Record ID: %s\n", d.Id())

	return resourceIpamSubnetsShow(d, meta)
}

func resourceIpamSubnetsShow(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	s := newSubnetFromResource(nil)

	subnet, err := client.Subnets().Show(d.Get("uuid").(string), s)
	log.Printf("HERE=> %+v\n\n", subnet)
	if err != nil {
		return fmt.Errorf("Could not find subnet: %s\n", err)
	}

	d.Set("uuid", subnet.ID)
	d.Set("name", subnet.Name)
	d.Set("tags", subnet.Tags)
	d.Set("metadata", subnet.Metadata)
	d.Set("pool", subnet.Pool)
	d.Set("start", subnet.Start)
	d.Set("end", subnet.End)

	return err
}

func resourceIpamSubnetsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	s := newSubnetFromResource(nil)

	if attr, ok := d.GetOk("uuid"); ok {
		s.ID = attr.(string)
	}
	if attr, ok := d.GetOk("name"); ok {
		s.Name = attr.(string)
	}
	if attr, ok := d.GetOk("tags"); ok {
		s.Tags = attr.([]string)
	}
	if attr, ok := d.GetOk("metadata"); ok {
		s.Metadata = attr
	}
	if attr, ok := d.GetOk("pool"); ok {
		s.Metadata = attr.(string)
	}
	if attr, ok := d.GetOk("start"); ok {
		s.Metadata = attr.(string)
	}
	if attr, ok := d.GetOk("end"); ok {
		s.Metadata = attr.(string)
	}

	log.Printf("Ipam Subnet update configuration: %#v", s)

	_, err := client.Subnets().UpdateShowSubnet(d.Get("uuid").(string), s)
	if err != nil {
		return fmt.Errorf("Failed to update Ipam Subnet configuration: %s\n", err)
	}

	return resourceIpamSubnetsShow(d, meta)
}

func resourceIpamSubnetsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	s := newSubnetFromResource(nil)

	log.Printf("Deleting Ipam Subnet: %s\n", d.Id())

	_, err := client.Subnets().Delete(d.Get("uuid").(string), s)
	if err != nil {
		return fmt.Errorf("Failed to delete Ipam Subnet: %s\n", err)
	}

	return nil
}

func newSubnetFromResource(d *schema.ResourceData) resources.SubnetV1 {
	subnet := resources.SubnetV1{}
	if d != nil {
		subnet.ID = d.Get("uuid").(string)
		subnet.Name = d.Get("name").(string)

		if v := d.Get("tags"); v != nil {
			for _, v := range v.([]interface{}) {
				subnet.Tags = append(subnet.Tags, v.(string))
			}
		}

		if v := d.Get("metadata"); v != nil {
			metadata := make(map[string]interface{})
			for k, v := range v.(map[string]interface{}) {
				metadata[k] = v
			}
			subnet.Metadata = metadata
		}

		subnet.Pool = d.Get("pool").(string)
		subnet.Start = d.Get("start").(string)
		subnet.End = d.Get("end").(string)
	}

	return subnet
}
