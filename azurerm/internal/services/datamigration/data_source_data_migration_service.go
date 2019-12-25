package datamigration

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmDataMigrationService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmDataMigrationServiceRead,

		Schema: map[string]*schema.Schema{
			"resource_group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocationForDataSource(),

			"virtual_subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sku_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"provisioning_state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmDataMigrationServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataMigration.ServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Data Migration Service (Service Name %q / Group Name %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading Data Migration Service (Service Name %q / Group Name %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("id", resp.ID)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("resource_group_name", resourceGroup)
	d.Set("name", resp.Name)
	if serviceProperties := resp.ServiceProperties; serviceProperties != nil {
		d.Set("provisioning_state", string(serviceProperties.ProvisioningState))
		d.Set("virtual_subnet_id", serviceProperties.VirtualSubnetID)
	}
	if err := d.Set("sku_name", resp.Sku.Name); err != nil {
		return fmt.Errorf("Error setting `sku_name`: %+v", err)
	}
	d.Set("type", resp.Type)
	d.Set("kind", resp.Kind)

	return tags.FlattenAndSet(d, resp.Tags)
}
