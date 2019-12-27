package mssqlvm

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmMsSqlVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmMsSqlVirtualMachineRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"virtual_machine_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sql_server_license_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sql_image_sku": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func dataSourceArmMsSqlVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQLVM.SQLVirtualMachinesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroupName := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resourceGroupName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q) was not found", name, resourceGroupName)
		}
		return fmt.Errorf("Error reading Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}

	d.SetId(*resp.ID)

	d.Set("resource_group", resourceGroupName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if properties := resp.Properties; properties != nil {
		d.Set("sql_image_sku", string(properties.SQLImageSku))
		d.Set("sql_server_license_type", string(properties.SQLServerLicenseType))
		d.Set("virtual_machine_resource_id", properties.VirtualMachineResourceID)
	}
	d.Set("name", name)
	d.Set("id", resp.ID)

	return tags.FlattenAndSet(d, resp.Tags)
}
