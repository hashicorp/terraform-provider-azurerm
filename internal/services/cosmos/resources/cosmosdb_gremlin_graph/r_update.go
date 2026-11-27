// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_gremlin_graph

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbGremlinGraphUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseGraphID(d.Id())
	if err != nil {
		return err
	}

	if err = common.CheckForChangeFromAutoscaleAndManualThroughput(d); err != nil {
		return fmt.Errorf("checking `autoscale_settings` and `throughput` for %s: %w", id, err)
	}

	partitionkeypaths := d.Get("partition_key_path").(string)

	db := cosmosdb.GremlinGraphCreateUpdateParameters{
		Properties: cosmosdb.GremlinGraphCreateUpdateProperties{
			Resource: cosmosdb.GremlinGraphResource{
				Id:             id.GraphName,
				IndexingPolicy: expandAzureRmCosmosDbGremlinGraphIndexingPolicy(d),
			},
			Options: &cosmosdb.CreateUpdateOptions{},
		},
	}

	if partitionkeypaths != "" {
		db.Properties.Resource.PartitionKey = &cosmosdb.ContainerPartitionKey{
			Paths: &[]string{partitionkeypaths},
			Kind:  pointer.To(cosmosdb.PartitionKindHash),
		}

		if partitionKeyVersion, ok := d.GetOk("partition_key_version"); ok {
			db.Properties.Resource.PartitionKey.Version = pointer.To(int64(partitionKeyVersion.(int)))
		}
	}

	if keys := expandAzureRmCosmosDbGremlinGraphUniqueKeys(d.Get("unique_key").(*pluginsdk.Set)); keys != nil {
		db.Properties.Resource.UniqueKeyPolicy = &cosmosdb.UniqueKeyPolicy{
			UniqueKeys: keys,
		}
	}

	if v, ok := d.GetOk("analytical_storage_ttl"); ok {
		db.Properties.Resource.AnalyticalStorageTtl = pointer.To(int64(v.(int)))
	}

	if defaultTTL, hasDefaultTTL := d.GetOk("default_ttl"); hasDefaultTTL {
		db.Properties.Resource.DefaultTtl = pointer.To(int64(defaultTTL.(int)))
	}

	if err = client.GremlinResourcesCreateUpdateGremlinGraphThenPoll(ctx, *id, db); err != nil {
		return fmt.Errorf("updating %q: %+v", id, err)
	}

	if common.HasThroughputChange(d) {
		if err := client.GremlinResourcesUpdateGremlinGraphThroughputThenPoll(ctx, *id, common.ExpandCosmosDBThroughputSettingsUpdateParameters(d)); err != nil {
			return fmt.Errorf("setting Throughput for %s: %+v - If the graph has not been created with an initial throughput, you cannot configure it later", id, err)
		}
	}

	return resourceCosmosDbGremlinGraphRead(d, meta)
}
