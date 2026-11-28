// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_container

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbSQLContainerUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseContainerID(d.Id())
	if err != nil {
		return err
	}

	if err = common.CheckForChangeFromAutoscaleAndManualThroughput(d); err != nil {
		return fmt.Errorf("checking `autoscale_settings` and `throughput` for %s: %w", id, err)
	}

	indexingPolicy := common.ExpandAzureRmCosmosDbIndexingPolicy(d)
	if err := common.ValidateAzureRmCosmosDbIndexingPolicy(indexingPolicy); err != nil {
		return fmt.Errorf("validating `indexing_policy`: %w", err)
	}

	db := cosmosdb.SqlContainerCreateUpdateParameters{
		Properties: cosmosdb.SqlContainerCreateUpdateProperties{
			Resource: cosmosdb.SqlContainerResource{
				Id:             id.ContainerName,
				IndexingPolicy: indexingPolicy,
			},
			Options: &cosmosdb.CreateUpdateOptions{},
		},
	}

	db.Properties.Resource.PartitionKey = &cosmosdb.ContainerPartitionKey{
		Kind: pointer.ToEnum[cosmosdb.PartitionKind](d.Get("partition_key_kind").(string)),
	}

	if v, ok := d.GetOk("partition_key_paths"); ok {
		db.Properties.Resource.PartitionKey.Paths = pluginsdk.ExpandStringSlice(v.([]any))
	}

	if partitionKeyVersion, ok := d.GetOk("partition_key_version"); ok {
		db.Properties.Resource.PartitionKey.Version = pointer.To(int64(partitionKeyVersion.(int)))
	}

	if keys := expandCosmosSQLContainerUniqueKeys(d.Get("unique_key").(*pluginsdk.Set)); keys != nil {
		db.Properties.Resource.UniqueKeyPolicy = &cosmosdb.UniqueKeyPolicy{
			UniqueKeys: keys,
		}
	}

	if analyticalStorageTTL, ok := d.GetOk("analytical_storage_ttl"); ok {
		db.Properties.Resource.AnalyticalStorageTtl = pointer.To(int64(analyticalStorageTTL.(int)))
	}

	if defaultTTL, hasTTL := d.GetOk("default_ttl"); hasTTL {
		db.Properties.Resource.DefaultTtl = pointer.To(int64(defaultTTL.(int)))
	}

	if err = client.SqlResourcesCreateUpdateSqlContainerThenPoll(ctx, *id, db); err != nil {
		return fmt.Errorf("updating %q: %+v", id, err)
	}

	if common.HasThroughputChange(d) {
		if err := client.SqlResourcesUpdateSqlContainerThroughputThenPoll(ctx, *id, common.ExpandCosmosDBThroughputSettingsUpdateParameters(d)); err != nil {
			return fmt.Errorf("setting Throughput for %s: %+v - If the collection has not been created with an initial throughput, you cannot configure it later", id, err)
		}
	}

	return resourceCosmosDbSQLContainerRead(d, meta)
}
