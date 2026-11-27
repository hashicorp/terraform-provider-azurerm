// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_gremlin_graph

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbGremlinGraphRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseGraphID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GremlinResourcesGetGremlinGraph(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.GraphName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("account_name", id.DatabaseAccountName)
	d.Set("database_name", id.GremlinDatabaseName)

	if model := resp.Model; model != nil {
		if graphProperties := model.Properties; graphProperties != nil {
			if props := graphProperties.Resource; props != nil {
				if pk := props.PartitionKey; pk != nil {
					if paths := pk.Paths; paths != nil {
						if l := len(*paths); l > 1 {
							return fmt.Errorf("retrieving `partition_key_path`, expected at most 1 path, got %d paths", l)
						}
						d.Set("partition_key_path", (*paths)[0])
					}

					if version := pk.Version; version != nil {
						d.Set("partition_key_version", version)
					}
				}

				if ip := props.IndexingPolicy; ip != nil {
					if err := d.Set("index_policy", flattenAzureRmCosmosDBGremlinGraphIndexingPolicy(props.IndexingPolicy)); err != nil {
						return fmt.Errorf("setting `index_policy`: %+v", err)
					}
				}

				if crp := props.ConflictResolutionPolicy; crp != nil {
					if err := d.Set("conflict_resolution_policy", common.FlattenCosmosDbConflictResolutionPolicy(crp)); err != nil {
						return fmt.Errorf("setting `conflict_resolution_policy`: %+v", err)
					}
				}

				if ukp := props.UniqueKeyPolicy; ukp != nil {
					if err := d.Set("unique_key", flattenCosmosGremlinGraphUniqueKeys(ukp.UniqueKeys)); err != nil {
						return fmt.Errorf("setting `unique_key`: %+v", err)
					}
				}

				if v := props.AnalyticalStorageTtl; v != nil {
					d.Set("analytical_storage_ttl", v)
				}

				if defaultTTL := props.DefaultTtl; defaultTTL != nil {
					d.Set("default_ttl", defaultTTL)
				}
			}
		}
	}

	databaseAccountID := cosmosdb.NewDatabaseAccountID(id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName)
	accResp, err := client.DatabaseAccountsGet(ctx, databaseAccountID)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", databaseAccountID, err)
	}

	if !helpers.IsServerlessCapacityMode(accResp.Model) {
		throughputResp, err := client.GremlinResourcesGetGremlinGraphThroughput(ctx, *id)
		if err != nil {
			if !response.WasNotFound(throughputResp.HttpResponse) {
				return fmt.Errorf("retrieving Throughput for %s: %v", id, err)
			} else {
				d.Set("throughput", nil)
				d.Set("autoscale_settings", nil)
			}
		} else {
			common.SetResourceDataThroughputFromResponse(*throughputResp.Model, d)
		}
	}
	return nil
}
