package dns

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2018-05-01/dns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dns/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceDnsZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDnsZoneRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"number_of_record_sets": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_number_of_record_sets": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"name_servers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceDnsZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.ZonesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	var (
		resp dns.Zone
		err  error
	)
	if resourceGroup != "" {
		resp, err = client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error: DNS Zone %q (Resource Group %q) was not found", name, resourceGroup)
			}
			return fmt.Errorf("Error reading DNS Zone %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	} else {
		var zone *dns.Zone
		zone, resourceGroup, err = findZone(client, ctx, name)
		if err != nil {
			return err
		}

		if zone == nil {
			return fmt.Errorf("Error: DNS Zone %q was not found", name)
		}

		resp = *zone
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("failed reading ID for DNS Zone %q (Resource Group %q)", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	if props := resp.ZoneProperties; props != nil {
		d.Set("number_of_record_sets", props.NumberOfRecordSets)
		d.Set("max_number_of_record_sets", props.MaxNumberOfRecordSets)

		nameServers := make([]string, 0)
		if ns := props.NameServers; ns != nil {
			nameServers = *ns
		}
		if err := d.Set("name_servers", nameServers); err != nil {
			return err
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func findZone(client *dns.ZonesClient, ctx context.Context, name string) (*dns.Zone, string, error) {
	zonesIterator, err := client.ListComplete(ctx, nil)
	if err != nil {
		return nil, "", fmt.Errorf("listing DNS Zones: %+v", err)
	}

	var found *dns.Zone
	for zonesIterator.NotDone() {
		zone := zonesIterator.Value()
		if zone.Name != nil && *zone.Name == name {
			if found != nil {
				return nil, "", fmt.Errorf("found multiple DNS zones with name %q, please specify the resource group", name)
			}
			found = &zone
		}
		if err := zonesIterator.NextWithContext(ctx); err != nil {
			return nil, "", fmt.Errorf("listing DNS Zones: %+v", err)
		}
	}

	if found == nil || found.ID == nil {
		return nil, "", fmt.Errorf("could not find DNS zone with name: %q", name)
	}

	id, err := parse.DnsZoneID(*found.ID)
	if err != nil {
		return nil, "", fmt.Errorf("DNS zone id not valid: %+v", err)
	}
	return found, id.ResourceGroup, nil
}
