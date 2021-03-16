package compute

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceDedicatedHost() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDedicatedHostRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDedicatedHostName(),
			},

			"dedicated_host_group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDedicatedHostGroupName(),
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceDedicatedHostRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)
	hostGroupName := d.Get("dedicated_host_group_name").(string)

	resp, err := client.Get(ctx, resourceGroupName, hostGroupName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Dedicated Host %q (Host Group Name %q / Resource Group %q) was not found", name, hostGroupName, resourceGroupName)
		}
		return fmt.Errorf("Error reading Dedicated Host %q (Host Group Name %q / Resource Group %q): %+v", name, hostGroupName, resourceGroupName, err)
	}

	d.SetId(*resp.ID)
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroupName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("dedicated_host_group_name", hostGroupName)

	return tags.FlattenAndSet(d, resp.Tags)
}
