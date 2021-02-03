package mssql

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/helper"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceMsSqlDatabase() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMsSqlDatabaseRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: helper.ValidateMsSqlDatabaseName,
			},

			"server_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ServerID,
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
				Type:     schema.TypeBool,
				Computed: true,
			},

			"sku_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"storage_account_type": {
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

func dataSourceMsSqlDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	mssqlServerId := d.Get("server_id").(string)
	serverId, err := parse.ServerID(mssqlServerId)
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, serverId.ResourceGroup, serverId.Name, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Database %q (Resource Group %q, SQL Server %q) was not found", name, serverId.ResourceGroup, serverId.Name)
		}

		return fmt.Errorf("making Read request on AzureRM Database %s (Resource Group %q, SQL Server %q): %+v", name, serverId.ResourceGroup, serverId.Name, err)
	}

	if id := resp.ID; id != nil {
		d.SetId(*resp.ID)
	}
	d.Set("name", name)
	d.Set("server_id", mssqlServerId)

	if props := resp.DatabaseProperties; props != nil {
		d.Set("collation", props.Collation)
		d.Set("elastic_pool_id", props.ElasticPoolID)
		d.Set("license_type", props.LicenseType)
		if props.MaxSizeBytes != nil {
			d.Set("max_size_gb", int32((*props.MaxSizeBytes)/int64(1073741824)))
		}
		d.Set("read_replica_count", props.ReadReplicaCount)
		if props.ReadScale == sql.DatabaseReadScaleEnabled {
			d.Set("read_scale", true)
		} else if props.ReadScale == sql.DatabaseReadScaleDisabled {
			d.Set("read_scale", false)
		}
		d.Set("sku_name", props.CurrentServiceObjectiveName)
		d.Set("storage_account_type", props.StorageAccountType)
		d.Set("zone_redundant", props.ZoneRedundant)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
