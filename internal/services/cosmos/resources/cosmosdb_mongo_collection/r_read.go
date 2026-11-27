// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_mongo_collection

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

func resourceCosmosDbMongoCollectionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseMongodbDatabaseCollectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.MongoDBResourcesGetMongoDBCollection(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("account_name", id.DatabaseAccountName)
	d.Set("database_name", id.MongodbDatabaseName)
	d.Set("name", id.CollectionName)

	databaseAccountID := cosmosdb.NewDatabaseAccountID(id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName)
	accResp, err := client.DatabaseAccountsGet(ctx, databaseAccountID)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", databaseAccountID, err)
	}

	if resp.Model != nil {
		if props := resp.Model.Properties; props != nil {
			if res := props.Resource; res != nil {
				// you can only have one
				if l := len(pointer.From(res.ShardKey)); l > 2 {
					return fmt.Errorf("unexpected number of shard keys: %d", l)
				}

				for k := range pointer.From(res.ShardKey) {
					d.Set("shard_key", k)
				}

				accountIsVersion36 := false
				if accResp.Model != nil {
					if accProps := accResp.Model.Properties; accProps != nil {
						if capabilities := accProps.Capabilities; capabilities != nil {
							for _, v := range *capabilities {
								if pointer.From(v.Name) == "EnableMongo" {
									accountIsVersion36 = true
								}
							}
						}
					}
				}

				indexes, systemIndexes, ttl := flattenCosmosMongoCollectionIndex(res.Indexes, accountIsVersion36)
				// In fact, the Azure API does not return `ExpireAfterSeconds` aka `default_ttl_seconds` when `default_ttl_seconds` is not set in tf the config.
				// When "default_ttl_seconds" is set to nil, it will be set to 0 in state file. 0 is invalid value for `default_ttl_seconds` and could not pass tf validation.
				// So when `default_ttl_seconds` is not set in tf config, we should not set the value of `default_ttl_seconds` but keep null in the state file.
				if ttl != nil {
					if err := d.Set("default_ttl_seconds", ttl); err != nil {
						return fmt.Errorf("setting `default_ttl_seconds`: %+v", err)
					}
				}
				if err := d.Set("index", indexes); err != nil {
					return fmt.Errorf("setting `index`: %+v", err)
				}
				if err := d.Set("system_indexes", systemIndexes); err != nil {
					return fmt.Errorf("setting `system_indexes`: %+v", err)
				}

				d.Set("analytical_storage_ttl", res.AnalyticalStorageTtl)
			}
		}
	}

	// if the cosmos account is serverless calling the get throughput api would yield an error
	if !helpers.IsServerlessCapacityMode(accResp.Model) {
		throughputResp, err := client.MongoDBResourcesGetMongoDBCollectionThroughput(ctx, *id)
		if err != nil {
			if !response.WasNotFound(throughputResp.HttpResponse) {
				return fmt.Errorf("retrieving Throughput for %s: %+v", id, err)
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

func flattenCosmosMongoCollectionIndex(input *[]cosmosdb.MongoIndex, accountIsVersion36 bool) (*[]map[string]any, *[]map[string]any, *int64) {
	indexes := make([]map[string]any, 0)
	systemIndexes := make([]map[string]any, 0)
	var ttl *int64
	if input == nil {
		return &indexes, &systemIndexes, ttl
	}

	for _, v := range *input {
		index := map[string]any{}
		systemIndex := map[string]any{}

		if v.Key != nil && v.Key.Keys != nil && len(*v.Key.Keys) > 0 {
			key := (*v.Key.Keys)[0]

			switch key {
			// As `DocumentDBDefaultIndex` and `_id` cannot be updated, so they would be moved into `system_indexes`.
			case "_id":
				systemIndex["keys"] = pluginsdk.FlattenSlice(v.Key.Keys)
				// The system index `_id` is always unique but api returns nil and it would be converted to `false` by zero-value. So it has to be manually set as `true`.
				systemIndex["unique"] = true

				systemIndexes = append(systemIndexes, systemIndex)

				if accountIsVersion36 {
					index["keys"] = pluginsdk.FlattenSlice(v.Key.Keys)
					index["unique"] = true
					indexes = append(indexes, index)
				}

			case "DocumentDBDefaultIndex":
				// Updating system index `DocumentDBDefaultIndex` is not a supported scenario.
				systemIndex["keys"] = pluginsdk.FlattenSlice(v.Key.Keys)

				isUnique := false
				if v.Options != nil && v.Options.Unique != nil {
					isUnique = *v.Options.Unique
				}
				systemIndex["unique"] = isUnique

				systemIndexes = append(systemIndexes, systemIndex)
			case "_ts":
				if v.Options != nil && v.Options.ExpireAfterSeconds != nil {
					// As `ExpireAfterSeconds` only can be applied to system index `_ts`, so it would be set in `default_ttl_seconds`.
					ttl = v.Options.ExpireAfterSeconds
				}
			default:
				// The other settable indexes would be set in `index`
				index["keys"] = pluginsdk.FlattenSlice(v.Key.Keys)

				isUnique := false
				if v.Options != nil && v.Options.Unique != nil {
					isUnique = *v.Options.Unique
				}
				index["unique"] = isUnique

				indexes = append(indexes, index)
			}
		}
	}

	return &indexes, &systemIndexes, ttl
}
