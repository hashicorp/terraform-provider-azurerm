// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_mongo_collection

import (
	"fmt"
	"strings"

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

func resourceCosmosDbMongoCollectionCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := cosmosdb.NewMongodbDatabaseCollectionID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("database_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.MongoDBResourcesGetMongoDBCollection(ctx, id)
		if !response.WasNotFound(existing.HttpResponse) {
			if err != nil {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
			return tf.ImportAsExistsError("azurerm_cosmosdb_mongo_collection", id.ID())
		}
	}

	var ttl *int
	if v, ok := d.GetOk("default_ttl_seconds"); ok {
		ttl = pointer.To(v.(int))
	}

	indexes, hasIdKey := expandCosmosMongoCollectionIndex(d.Get("index").(*pluginsdk.Set).List(), ttl)
	if !hasIdKey {
		return fmt.Errorf("index with '_id' key is required")
	}

	db := cosmosdb.MongoDBCollectionCreateUpdateParameters{
		Properties: cosmosdb.MongoDBCollectionCreateUpdateProperties{
			Resource: cosmosdb.MongoDBCollectionResource{
				Id:      id.CollectionName,
				Indexes: indexes,
			},
			Options: &cosmosdb.CreateUpdateOptions{},
		},
	}

	if analyticalStorageTTL, ok := d.GetOk("analytical_storage_ttl"); ok {
		db.Properties.Resource.AnalyticalStorageTtl = pointer.To(int64(analyticalStorageTTL.(int)))
	}

	if throughput, hasThroughput := d.GetOk("throughput"); hasThroughput {
		if throughput != 0 {
			db.Properties.Options.Throughput = common.ConvertThroughputFromResourceData(throughput)
		}
	}

	if _, hasAutoscaleSettings := d.GetOk("autoscale_settings"); hasAutoscaleSettings {
		db.Properties.Options.AutoScaleSettings = common.ExpandCosmosDbAutoscaleSettings(d)
	}

	if shardKey := d.Get("shard_key").(string); shardKey != "" {
		db.Properties.Resource.ShardKey = pointer.To(map[string]string{
			shardKey: "Hash", // looks like only hash is supported for now
		})
	}

	if err := client.MongoDBResourcesCreateUpdateMongoDBCollectionCallbackThenPoll(ctx, id, db, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosDbMongoCollectionRead(d, meta)
}

func expandCosmosMongoCollectionIndex(indexes []any, defaultTtl *int) (*[]cosmosdb.MongoIndex, bool) {
	results := make([]cosmosdb.MongoIndex, 0)

	hasIdKey := false

	if len(indexes) != 0 {
		for _, v := range indexes {
			index := v.(map[string]any)
			keys := index["keys"].([]any)

			for _, key := range keys {
				if strings.EqualFold("_id", key.(string)) {
					hasIdKey = true
				}
			}

			results = append(results, cosmosdb.MongoIndex{
				Key: &cosmosdb.MongoIndexKeys{
					Keys: pluginsdk.ExpandStringSlice(index["keys"].([]any)),
				},
				Options: &cosmosdb.MongoIndexOptions{
					Unique: pointer.To(index["unique"].(bool)),
				},
			})
		}
	}

	if defaultTtl != nil {
		results = append(results, cosmosdb.MongoIndex{
			Key: &cosmosdb.MongoIndexKeys{
				Keys: &[]string{"_ts"},
			},
			Options: &cosmosdb.MongoIndexOptions{
				ExpireAfterSeconds: pointer.To(int64(*defaultTtl)),
			},
		})
	}

	return &results, hasIdKey
}
