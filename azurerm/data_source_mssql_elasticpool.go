package azurerm

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmMsSqlElasticpool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmMsSqlElasticpoolRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"server_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"location": azure.SchemaLocationForDataSource(),

			"max_size_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_size_gb": {
				Type:     schema.TypeFloat,
				Computed: true,
			},

			"per_db_min_capacity": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"per_db_max_capacity": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),

			"zone_redundant": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceArmMsSqlElasticpoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).MSSQL.ElasticPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	elasticPoolName := d.Get("name").(string)
	serverName := d.Get("server_name").(string)

	resp, err := client.Get(ctx, resourceGroup, serverName, elasticPoolName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Elasticpool %q (Resource Group %q, SQL Server %q) was not found", elasticPoolName, resourceGroup, serverName)
		}

		return fmt.Errorf("Error making Read request on AzureRM Elasticpool %s (Resource Group %q, SQL Server %q): %+v", elasticPoolName, resourceGroup, serverName, err)
	}

	if id := resp.ID; id != nil {
		d.SetId(*resp.ID)
	}
	d.Set("name", elasticPoolName)
	d.Set("resource_group_name", resourceGroup)
	d.Set("server_name", serverName)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.ElasticPoolProperties; props != nil {
		d.Set("max_size_gb", float64(*props.MaxSizeBytes/int64(1073741824)))
		d.Set("max_size_bytes", props.MaxSizeBytes)

		d.Set("zone_redundant", props.ZoneRedundant)

		if perDbSettings := props.PerDatabaseSettings; perDbSettings != nil {
			d.Set("per_db_min_capacity", perDbSettings.MinCapacity)
			d.Set("per_db_max_capacity", perDbSettings.MaxCapacity)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
