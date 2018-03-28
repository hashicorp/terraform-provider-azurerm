package azurerm

import (
	"fmt"

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

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

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

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading DNS Zone %q (Resource Group %q): %+v", name, resourceGroup, err)
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
