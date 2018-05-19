package azurerm

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2016-04-01/dns"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmDnsZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmDnsZoneRead,

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
				Type:     schema.TypeString,
				Computed: true,
			},

			"max_number_of_record_sets": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name_servers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmDnsZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).zonesClient
	ctx := meta.(*ArmClient).StopContext

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
		rgClient := meta.(*ArmClient).resourceGroupsClient

		resp, resourceGroup, err = findZone(client, rgClient, ctx, name)
		if err != nil {
			return err
		}

		if resourceGroup == "" {
			return fmt.Errorf("Error: DNS Zone %q was not found", name)
		}
	}

	d.SetId(*resp.ID)
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	if props := resp.ZoneProperties; props != nil {
		d.Set("number_of_record_sets", props.NumberOfRecordSets)
		d.Set("max_number_of_record_sets", props.MaxNumberOfRecordSets)

		if ns := props.NameServers; ns != nil {
			nameServers := make([]string, 0, len(*ns))
			for _, ns := range *ns {
				nameServers = append(nameServers, ns)
			}
			if err := d.Set("name_servers", nameServers); err != nil {
				return err
			}
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func findZone(client dns.ZonesClient, rgClient resources.GroupsClient, ctx context.Context, name string) (dns.Zone, string, error) {
	groups, err := rgClient.List(ctx, "", nil)
	if err != nil {
		return dns.Zone{}, "", fmt.Errorf("Error listing Resource Groups: %+v", err)
	}

	for _, g := range groups.Values() {
		resourceGroup := *g.Name

		zones, err := client.ListByResourceGroup(ctx, resourceGroup, nil)
		if err != nil {
			return dns.Zone{}, "", fmt.Errorf("Error listing DNS Zones (Resource Group: %s): %+v", resourceGroup, err)
		}

		for _, z := range zones.Values() {
			if *z.Name == name {
				return z, resourceGroup, nil
			}
		}
	}

	return dns.Zone{}, "", nil
}
