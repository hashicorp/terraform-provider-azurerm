package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceServiceTags() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServiceTagsRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"location": azure.SchemaLocation(),

			"service": {
				Type:     schema.TypeString,
				Required: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"address_prefixes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceServiceTagsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ServiceTagsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	location := d.Get("location").(string)
	res, err := client.List(ctx, location)
	if err != nil {
		return err
	}

	service := d.Get("service").(string)
	region := d.Get("region").(string)

	if res.Values == nil {
		d.Set("address_prefixes", make([]string, 0))
		return nil
	}

	addressPrefixes := make([]string, 0)
	for _, sti := range *res.Values {
		if props := sti.Properties; props != nil {
			if *props.SystemService != service {
				continue
			}

			if *props.Region == region {
				addressPrefixes = *props.AddressPrefixes
			}
		}
	}

	err = d.Set("address_prefixes", addressPrefixes)
	if err != nil {
		return fmt.Errorf("Error setting `address_prefixes`: %+v", err)
	}

	id := "servicetags-" + location + "-" + service
	if region != "" {
		id = id + "-" + region
	}
	d.SetId(id)

	return nil
}
