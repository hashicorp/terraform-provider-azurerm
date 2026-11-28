// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_container

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbSQLContainerRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseContainerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.SqlResourcesGetSqlContainer(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ContainerName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("account_name", id.DatabaseAccountName)
	d.Set("database_name", id.SqlDatabaseName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if res := props.Resource; res != nil {
				if pk := res.PartitionKey; pk != nil {
					d.Set("partition_key_kind", pointer.FromEnum(pk.Kind))
					d.Set("partition_key_paths", pluginsdk.FlattenSlice(pk.Paths))
					d.Set("partition_key_version", pk.Version)
				}

				if ukp := res.UniqueKeyPolicy; ukp != nil {
					if err := d.Set("unique_key", flattenCosmosSQLContainerUniqueKeys(ukp.UniqueKeys)); err != nil {
						return fmt.Errorf("setting `unique_key`: %+v", err)
					}
				}

				d.Set("analytical_storage_ttl", res.AnalyticalStorageTtl)
				d.Set("default_ttl", res.DefaultTtl)
				d.Set("indexing_policy", common.FlattenAzureRmCosmosDbIndexingPolicy(res.IndexingPolicy))

				if err := d.Set("conflict_resolution_policy", common.FlattenCosmosDbConflictResolutionPolicy(res.ConflictResolutionPolicy)); err != nil {
					return fmt.Errorf("setting `conflict_resolution_policy`: %+v", err)
				}
			}
		}
	}

	databaseAccountID := cosmosdb.NewDatabaseAccountID(id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName)
	accResp, err := client.DatabaseAccountsGet(ctx, databaseAccountID)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", databaseAccountID, err)
	}

	// if the cosmos account is serverless calling the get throughput api would yield an error
	if !helpers.IsServerlessCapacityMode(accResp.Model) {
		throughputResp, err := client.SqlResourcesGetSqlContainerThroughput(ctx, *id)
		if err != nil {
			if !response.WasNotFound(throughputResp.HttpResponse) {
				return fmt.Errorf("retrieving Throughput for %s: %v", id, err)
			} else {
				d.Set("throughput", nil)
				d.Set("autoscale_settings", nil)
			}
		} else {
			common.SetResourceDataThroughputFromResponse(pointer.From(throughputResp.Model), d)
		}
	}

	return nil
}

func flattenCosmosSQLContainerUniqueKeys(keys *[]cosmosdb.UniqueKey) *[]map[string]any {
	if keys == nil {
		return nil
	}

	slice := make([]map[string]any, 0)
	for _, k := range *keys {
		if k.Paths == nil {
			continue
		}

		slice = append(slice, map[string]any{
			"paths": *k.Paths,
		})
	}

	return &slice
}
