package mssql

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmMsSqlDatabase() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmMsSqlDatabaseRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateMsSqlDatabaseName,
			},

			"sql_server_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ValidateMsSqlServerID,
			},

			"collation": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"elastic_pool_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"license_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"max_size_gb": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"read_replica_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"read_scale": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sku_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"zone_redundant": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmMsSqlDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	mssqlServerId := d.Get("sql_server_id").(string)
	serverId, err := parse.MsSqlServerID(mssqlServerId)
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, serverId.ResourceGroup, serverId.Name, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Database %q (Resource Group %q, SQL Server %q) was not found", name, serverId.ResourceGroup, serverId.Name)
		}

		return fmt.Errorf("Error making Read request on AzureRM Database %s (Resource Group %q, SQL Server %q): %+v", name, serverId.ResourceGroup, serverId.Name, err)
	}

	if id := resp.ID; id != nil {
		d.SetId(*resp.ID)
	}
	d.Set("name", name)
	d.Set("sql_server_id", mssqlServerId)

	if props := resp.DatabaseProperties; props != nil {
		d.Set("collation", props.Collation)

		if props.ElasticPoolID != nil {
			d.Set("elastic_pool_id", props.ElasticPoolID)
		}

		d.Set("license_type", props.LicenseType)
		d.Set("max_size_gb", int32((*props.MaxSizeBytes)/int64(1073741824)))

		if props.ReadReplicaCount != nil {
			d.Set("read_replica_count", props.ReadReplicaCount)
		}
		d.Set("read_scale", props.ReadScale)
		d.Set("sku_name", props.CurrentServiceObjectiveName)
		d.Set("zone_redundant", props.ZoneRedundant)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
