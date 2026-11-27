// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_gremlin_graph

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

func resourceCosmosDbGremlinGraphCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := cosmosdb.NewGraphID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("database_name").(string), d.Get("name").(string))
	partitionkeypaths := d.Get("partition_key_path").(string)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.GremlinResourcesGetGremlinGraph(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_cosmosdb_gremlin_graph", id.ID())
		}
	}

	db := cosmosdb.GremlinGraphCreateUpdateParameters{
		Properties: cosmosdb.GremlinGraphCreateUpdateProperties{
			Resource: cosmosdb.GremlinGraphResource{
				Id:                       id.GraphName,
				IndexingPolicy:           expandAzureRmCosmosDbGremlinGraphIndexingPolicy(d),
				ConflictResolutionPolicy: common.ExpandCosmosDbConflicResolutionPolicy(d.Get("conflict_resolution_policy").([]any)),
			},
			Options: &cosmosdb.CreateUpdateOptions{},
		},
	}

	if v, ok := d.GetOk("analytical_storage_ttl"); ok {
		db.Properties.Resource.AnalyticalStorageTtl = pointer.To(int64(v.(int)))
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

	if defaultTTL, hasDefaultTTL := d.GetOk("default_ttl"); hasDefaultTTL {
		if defaultTTL != 0 {
			db.Properties.Resource.DefaultTtl = pointer.To(int64(defaultTTL.(int)))
		}
	}

	if throughput, hasThroughput := d.GetOk("throughput"); hasThroughput {
		if throughput != 0 {
			db.Properties.Options.Throughput = common.ConvertThroughputFromResourceData(throughput)
		}
	}

	if _, hasAutoscaleSettings := d.GetOk("autoscale_settings"); hasAutoscaleSettings {
		db.Properties.Options.AutoScaleSettings = common.ExpandCosmosDbAutoscaleSettings(d)
	}

	if err := client.GremlinResourcesCreateUpdateGremlinGraphCallbackThenPoll(ctx, id, db, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosDbGremlinGraphRead(d, meta)
}
