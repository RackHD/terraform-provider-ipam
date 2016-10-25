package main

import (
	"fmt"
	"log"

	"github.com/RackHD/ipam/resources"
	"github.com/ekaley/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIpamPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpamPoolsCreate,
		Read:   resourceIpamPoolsShow,
		Update: resourceIpamPoolsUpdate,
		Delete: resourceIpamPoolsDelete,

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
		},
	}
}

func resourceIpamPoolsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	p := newPoolFromResource(d)

	pool, err := client.Pools().CreateShowPool(p)
	if err != nil {
		return fmt.Errorf("Failed to create Pool: %s\n", err)
	}
	d.SetId(pool.ID)
	d.Set("uuid", pool.ID)
	log.Printf("Record ID: %s\n", d.Id())

	return resourceIpamPoolsShow(d, meta)
}

func resourceIpamPoolsShow(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	p := newPoolFromResource(nil)

	pool, err := client.Pools().Show(d.Get("uuid").(string), p)
	if err != nil {
		return fmt.Errorf("Could not find Pool: %s\n", err)
	}

	d.Set("uuid", pool.ID)
	d.Set("name", pool.Name)
	d.Set("tags", pool.Tags)
	d.Set("metadata", pool.Metadata)

	return err
}

func resourceIpamPoolsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	p := newPoolFromResource(nil)

	if attr, ok := d.GetOk("uuid"); ok {
		p.ID = attr.(string)
	}
	if attr, ok := d.GetOk("name"); ok {
		p.Name = attr.(string)
	}
	if attr, ok := d.GetOk("tags"); ok {
		p.Tags = attr.([]string)
	}
	if attr, ok := d.GetOk("metadata"); ok {
		p.Metadata = attr
	}

	log.Printf("Ipam pool update configuration: %#v", p)

	_, err := client.Pools().UpdateShowPool(d.Get("uuid").(string), p)
	if err != nil {
		return fmt.Errorf("Failed to update Ipam Pool configuration: %s\n", err)
	}

	return resourceIpamPoolsShow(d, meta)
}

func resourceIpamPoolsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	p := newPoolFromResource(nil)

	log.Printf("Deleting Ipam Pool: %s\n", d.Id())

	_, err := client.Pools().Delete(d.Get("uuid").(string), p)
	if err != nil {
		return fmt.Errorf("Failed to delete Ipam Pool: %s\n", err)
	}

	return nil
}

func newPoolFromResource(d *schema.ResourceData) resources.PoolV1 {
	pool := resources.PoolV1{}
	if d != nil {
		pool.ID = d.Get("uuid").(string)
		pool.Name = d.Get("name").(string)

		if v := d.Get("tags"); v != nil {
			for _, v := range v.([]interface{}) {
				pool.Tags = append(pool.Tags, v.(string))
			}
		}

		if v := d.Get("metadata"); v != nil {
			metadata := make(map[string]interface{})
			for k, v := range v.(map[string]interface{}) {
				metadata[k] = v
			}
			pool.Metadata = metadata
		}
	}

	return pool
}
