package sql

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceSqlDatabase() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSqlDatabaseRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

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
			"blob_extended_auditing_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"retention_days": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"audit_actions_and_groups": {
							Type:     schema.TypeSet,
							Computed: true,
							Set:      schema.HashString,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"storage_account_subscription_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_storage_secondary_key_in_use": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"predicate_expression": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceArmSqlDatabaseRead(d *schema.ResourceData, meta interface{}) error {
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

	auditingClient := meta.(*clients.Client).Sql.ExtendedDatabaseBlobAuditingPoliciesClient
	auditingResp, err := auditingClient.Get(ctx, resourceGroup, serverName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading SQL Server %q Database %q Blob Auditing Policies - removing from state", serverName, name)
		}
		return fmt.Errorf("Error reading SQL Server %s Database %q: %v Blob Auditing Policies", serverName, name, err)
	}

	d.Set("blob_extended_auditing_policy", flattenAzureRmSqlDBBlobAuditingPolicies(&auditingResp, d))

	return tags.FlattenAndSet(d, resp.Tags)
}
