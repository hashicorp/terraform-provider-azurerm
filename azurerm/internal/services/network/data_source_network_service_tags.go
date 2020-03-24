package network

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceNetworkServiceTags() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkServiceTagsRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"location": azure.SchemaLocation(),

			"service": {
				Type:     schema.TypeString,
				Required: true,
			},

			"location_filter": {
				Type:             schema.TypeString,
				Optional:         true,
				StateFunc:        azure.NormalizeLocation,
				DiffSuppressFunc: azure.SuppressLocationDiff,
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

func dataSourceNetworkServiceTagsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ServiceTagsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	location := azure.NormalizeLocation(d.Get("location").(string))
	res, err := client.List(ctx, location)
	if err != nil {
		return err
	}

	if res.Values == nil {
		return errors.New("unexpected nil value for service tag information")
	}

	service := d.Get("service").(string)
	locationFilter := d.Get("location_filter").(string)

	for _, sti := range *res.Values {
		if sti.Name == nil || !isServiceTagOf(*sti.Name, service) {
			continue
		}

		if props := sti.Properties; props != nil {
			if *props.Region == locationFilter {
				addressPrefixes := make([]string, 0)
				if props.AddressPrefixes != nil {
					addressPrefixes = *props.AddressPrefixes
				}
				err = d.Set("address_prefixes", addressPrefixes)
				if err != nil {
					return fmt.Errorf("Error setting `address_prefixes`: %+v", err)
				}

				if sti.ID == nil {
					return errors.New("unexpected nil ID for service tag")
				}

				d.SetId(*sti.ID)
				return nil
			}
		}
	}
	return errors.New("specified service tag not found")
}

// isServiceTagOf is used to check whether a service tag name belongs to the service of name `serviceName`.
// Service tag name has format as below:
// - (regional) serviceName.locationName
// - (all) serviceName
func isServiceTagOf(stName, serviceName string) bool {
	stNameComponents := strings.Split(stName, ".")
	if len(stNameComponents) != 1 && len(stNameComponents) != 2 {
		return false
	}
	return stNameComponents[0] == serviceName
}
