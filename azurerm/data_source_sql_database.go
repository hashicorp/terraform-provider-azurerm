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
			"collation": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_secondary_location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"edition": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"elastic_pool_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"failover_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"read_scale": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"requested_service_objective_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"requested_service_objective_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"server_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func dataSourceArmSqlDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sql.DatabasesClient
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

	if id := resp.ID; id != nil {
		d.SetId(*id)
	}

	d.Set("collation", resp.Collation)

	if dsLocation := resp.ElasticPoolName; dsLocation != nil {
		d.Set("default_secondary_location", dsLocation)
	}

	d.Set("edition", string(resp.Edition))

	if ep := resp.ElasticPoolName; ep != nil {
		d.Set("elastic_pool_name", ep)
	}

	if fogID := resp.FailoverGroupID; fogID != nil {
		d.Set("failover_group_id", fogID)
	}

	d.Set("location", azure.NormalizeLocation(*resp.Location))

	readScale := false
	if resp.ReadScale == sql.ReadScaleEnabled {
		readScale = true
	}
	d.Set("read_scale", readScale)

	if rsoID := resp.RequestedServiceObjectiveID; rsoID != nil {
		d.Set("requested_service_objective_id", rsoID)
	}

	if rsoName := resp.RequestedServiceObjectiveID; rsoName != nil {
		d.Set("requested_service_objective_name", rsoName)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
