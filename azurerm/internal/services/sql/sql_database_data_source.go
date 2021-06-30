package sql

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceSqlDatabase() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmSqlDatabaseRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"location": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"server_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"edition": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"collation": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"elastic_pool_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"default_secondary_location": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"read_scale": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"tags": tags.Schema(),

			"failover_group_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmSqlDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.DatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

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
