package cosmos

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-06-15/documentdb"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceCassandraManagedInstance() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceCassandraManagedInstanceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"cluster_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"tags": tags.SchemaDataSource(),

			"kind": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"location": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"delegated_management_subnet_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"initial_cassandra_admin_password": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCassandraManagedInstanceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraMIClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	clusterName := d.Get("cluster_name").(string)

	resp, err := client.GetCluster(ctx, resourceGroup, clusterName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Cluster name %q (Resource Group %q) was not found", clusterName, resourceGroup)
		}

		return fmt.Errorf("making Read request on AzureRM CosmosDB Account %s (Resource Group %q): %+v", clusterName, resourceGroup, err)
	}

	d.SetId(*resp.ID)
	d.Set("name", clusterName)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("kind", string(*resp.ID))

	if props := resp.CassandraKeyspaceGetProperties; props != nil {
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenAzureRmCassandraManagedInstanceCapabilitiesAsList(capabilities *[]documentdb.Capability) *[]map[string]interface{} {
	slice := make([]map[string]interface{}, 0)

	for _, c := range *capabilities {
		if v := c.Name; v != nil {
			e := map[string]interface{}{
				"name": *v,
			}
			slice = append(slice, e)
		}
	}

	return &slice
}

func flattenAzureRmCassandraManagedInstanceVirtualNetworkRulesAsList(rules *[]documentdb.VirtualNetworkRule) []map[string]interface{} {
	if rules == nil {
		return []map[string]interface{}{}
	}

	virtualNetworkRules := make([]map[string]interface{}, len(*rules))
	for i, r := range *rules {
		virtualNetworkRules[i] = map[string]interface{}{
			"id": *r.ID,
		}
	}
	return virtualNetworkRules
}
