// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_mongo_collection

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbMongoCollectionUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseMongodbDatabaseCollectionID(d.Id())
	if err != nil {
		return err
	}

	if err := common.CheckForChangeFromAutoscaleAndManualThroughput(d); err != nil {
		return fmt.Errorf("checking `autoscale_settings` and `throughput` for %s: %w", id, err)
	}

	existing, err := client.MongoDBResourcesGetMongoDBCollection(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %w", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: properties was nil", id)
	}

	if existing.Model.Properties.Resource == nil {
		return fmt.Errorf("retrieving %s: resource was nil", id)
	}

	db := cosmosdb.MongoDBCollectionCreateUpdateParameters{
		Properties: cosmosdb.MongoDBCollectionCreateUpdateProperties{
			Resource: cosmosdb.MongoDBCollectionResource{
				Id:       id.CollectionName,
				Indexes:  existing.Model.Properties.Resource.Indexes,
				ShardKey: existing.Model.Properties.Resource.ShardKey,
			},
			Options: &cosmosdb.CreateUpdateOptions{},
		},
	}

	if d.HasChanges("default_ttl_seconds", "index") {
		var ttl *int
		if v, ok := d.GetOk("default_ttl_seconds"); ok {
			ttl = pointer.To(v.(int))
		}

		indexes, hasIdKey := expandCosmosMongoCollectionIndex(d.Get("index").(*pluginsdk.Set).List(), ttl)
		if !hasIdKey {
			return errors.New("index with '_id' key is required")
		}
		db.Properties.Resource.Indexes = indexes
	}

	if d.HasChange("analytical_storage_ttl") {
		if analyticalStorageTTL, ok := d.GetOk("analytical_storage_ttl"); ok {
			db.Properties.Resource.AnalyticalStorageTtl = pointer.To(int64(analyticalStorageTTL.(int)))
		}
	}

	if err := client.MongoDBResourcesCreateUpdateMongoDBCollectionThenPoll(ctx, *id, db); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if common.HasThroughputChange(d) {
		if err := client.MongoDBResourcesUpdateMongoDBCollectionThroughputThenPoll(ctx, *id, common.ExpandCosmosDBThroughputSettingsUpdateParameters(d)); err != nil {
			return fmt.Errorf("setting Throughput for %s: %+v - if the collection has not been created with an initial throughput, you cannot configure it later", id, err)
		}
	}

	return resourceCosmosDbMongoCollectionRead(d, meta)
}
