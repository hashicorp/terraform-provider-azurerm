package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2015-05-01-preview/sql"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceSqlDatabase() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSqlDatabaseRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"server_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"edition": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"collation": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"elastic_pool_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_secondary_location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"read_scale": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"tags": tags.Schema(),

			"failover_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmSqlDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Sql.DatabasesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)

	resp, err := client.Get(ctx, resourceGroup, serverName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("SQL Database %q (server %q / resource group %q) was not found", name, serverName, resourceGroup)
		}

		return fmt.Errorf("Error retrieving SQL Database %q (server %q / resource group %q): %s", name, serverName, resourceGroup, err)
	}

	d.Set("location", azure.NormalizeLocation(*resp.Location))

	if id := resp.ID; id != nil {
		d.SetId(*id)
	}

	if props := resp.DatabaseProperties; props != nil {

		d.Set("collation", props.Collation)

		d.Set("default_secondary_location", props.DefaultSecondaryLocation)

		d.Set("edition", string(props.Edition))

		d.Set("elastic_pool_name", props.ElasticPoolName)

		d.Set("failover_group_id", props.FailoverGroupID)

		d.Set("read_scale", props.ReadScale == sql.ReadScaleEnabled)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
