package main

import (
	"fmt"
	"log"

	"github.com/RackHD/ipam/resources"
	"github.com/ekaley/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIpamReservation() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpamReservationsCreate,
		Read:   resourceIpamReservationsShow,
		Update: resourceIpamReservationsUpdate,
		Delete: resourceIpamReservationsDelete,

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
			"subnet": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceIpamReservationsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	r := newReservationFromResource(d)

	reservation, err := client.Reservations().CreateShowReservation(d.Get("subnet").(string), r)
	if err != nil {
		return fmt.Errorf("Failed to create Reservation: %s\n", err)
	}
	d.Set("uuid", reservation.ID)
	d.SetId(reservation.ID)
	log.Printf("Record ID: %s\n", d.Id())
	return resourceIpamReservationsShow(d, meta)
}

func resourceIpamReservationsShow(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	r := newReservationFromResource(nil)
	log.Printf("UUID: %s\n", d.Get("uuid").(string))
	log.Printf("SEND: %s\n", r)
	reservation, err := client.Reservations().Show(d.Get("uuid").(string), r)
	if err != nil {
		return fmt.Errorf("Could not find Reservation: %s\n", err)
	}
	lease, err := client.Leases().Index(reservation.ID)
	if err != nil {
		return fmt.Errorf("Could not find any Leases: %s\n", err)
	}

	d.Set("uuid", reservation.ID)
	d.Set("name", reservation.Name)
	d.Set("tags", reservation.Tags)
	d.Set("metadata", reservation.Metadata)
	d.Set("subnet", reservation.Subnet)
	d.Set("ip", lease.Leases[0].Address)

	return err
}

func resourceIpamReservationsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	r := newReservationFromResource(nil)

	if attr, ok := d.GetOk("uuid"); ok {
		r.ID = attr.(string)
	}
	if attr, ok := d.GetOk("name"); ok {
		r.Name = attr.(string)
	}
	if attr, ok := d.GetOk("tags"); ok {
		r.Tags = attr.([]string)
	}
	if attr, ok := d.GetOk("metadata"); ok {
		r.Metadata = attr
	}
	if attr, ok := d.GetOk("subnet"); ok {
		r.Metadata = attr.(string)
	}

	log.Printf("Ipam Reservation update configuration: %#v", r)

	_, err := client.Reservations().UpdateShowReservation(d.Get("uuid").(string), r)
	if err != nil {
		return fmt.Errorf("Failed to update Ipam Reservation configuration: %s\n", err)
	}

	return resourceIpamReservationsShow(d, meta)
}

func resourceIpamReservationsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	r := newReservationFromResource(nil)

	log.Printf("Deleting Ipam Reservation: %s\n", d.Get("uuid").(string))

	_, err := client.Reservations().Delete(d.Id(), r)
	if err != nil {
		return fmt.Errorf("Failed to delete Ipam Reservation: %s\n", err)
	}

	return nil
}

func newReservationFromResource(d *schema.ResourceData) resources.ReservationV1 {
	reservation := resources.ReservationV1{}
	if d != nil {
		reservation.ID = d.Get("uuid").(string)
		reservation.Name = d.Get("name").(string)

		if v := d.Get("tags"); v != nil {
			for _, v := range v.([]interface{}) {
				reservation.Tags = append(reservation.Tags, v.(string))
			}
		}

		if v := d.Get("metadata"); v != nil {
			metadata := make(map[string]interface{})
			for k, v := range v.(map[string]interface{}) {
				metadata[k] = v
			}
			reservation.Metadata = metadata
		}

		reservation.Subnet = d.Get("subnet").(string)
	}

	return reservation
}
