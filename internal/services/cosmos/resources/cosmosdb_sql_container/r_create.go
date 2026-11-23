// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_container

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbSQLContainerCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := cosmosdb.NewContainerID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("database_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.SqlResourcesGetSqlContainer(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_cosmosdb_sql_container", id.ID())
		}
	}

	indexingPolicy := common.ExpandAzureRmCosmosDbIndexingPolicy(d)
	if err := common.ValidateAzureRmCosmosDbIndexingPolicy(indexingPolicy); err != nil {
		return fmt.Errorf("validating `indexing_policy`: %w", err)
	}

	db := cosmosdb.SqlContainerCreateUpdateParameters{
		Properties: cosmosdb.SqlContainerCreateUpdateProperties{
			Resource: cosmosdb.SqlContainerResource{
				Id:                       id.ContainerName,
				IndexingPolicy:           indexingPolicy,
				ConflictResolutionPolicy: common.ExpandCosmosDbConflicResolutionPolicy(d.Get("conflict_resolution_policy").([]any)),
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

	if throughput, hasThroughput := d.GetOk("throughput"); hasThroughput {
		if throughput != 0 {
			db.Properties.Options.Throughput = common.ConvertThroughputFromResourceData(throughput)
		}
	}

	if _, hasAutoscaleSettings := d.GetOk("autoscale_settings"); hasAutoscaleSettings {
		db.Properties.Options.AutoScaleSettings = common.ExpandCosmosDbAutoscaleSettings(d)
	}

	if err := client.SqlResourcesCreateUpdateSqlContainerCallbackThenPoll(ctx, id, db, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosDbSQLContainerRead(d, meta)
}

func expandCosmosSQLContainerUniqueKeys(s *pluginsdk.Set) *[]cosmosdb.UniqueKey {
	i := s.List()
	if len(i) == 0 || i[0] == nil {
		return nil
	}

	keys := make([]cosmosdb.UniqueKey, 0)
	for _, k := range i {
		key := k.(map[string]any)

		paths := key["paths"].(*pluginsdk.Set).List()
		if len(paths) == 0 {
			continue
		}

		keys = append(keys, cosmosdb.UniqueKey{
			Paths: pluginsdk.ExpandStringSlice(paths),
		})
	}

	return &keys
}
