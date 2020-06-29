package databasemigration

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

func dataSourceArmDatabaseMigrationService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmDatabaseMigrationServiceRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDatabaseMigrationServiceName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sku_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmDatabaseMigrationServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DatabaseMigration.ServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Database Migration Service (Service Name %q / Group Name %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading Database Migration Service (Service Name %q / Group Name %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("unexpected empty ID retrieved for Database Migration Service (Service Name %q / Group Name %q)", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("resource_group_name", resourceGroup)
	d.Set("name", resp.Name)
	if serviceProperties := resp.ServiceProperties; serviceProperties != nil {
		d.Set("subnet_id", serviceProperties.VirtualSubnetID)
	}
	if resp.Sku != nil {
		d.Set("sku_name", resp.Sku.Name)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
