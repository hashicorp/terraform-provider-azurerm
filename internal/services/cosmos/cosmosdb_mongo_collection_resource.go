// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCosmosDbMongoCollection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCosmosDbMongoCollectionCreate,
		Read:   resourceCosmosDbMongoCollectionRead,
		Update: resourceCosmosDbMongoCollectionUpdate,
		Delete: resourceCosmosDbMongoCollectionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := cosmosdb.ParseMongodbDatabaseCollectionID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.MongoCollectionV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosEntityName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosAccountName,
			},

			"database_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosEntityName,
			},

			// SDK/api accepts an array.. but only one is allowed
			"shard_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// default TTL is simply an index on _ts with expireAfterOption, given we can't seem to set TTLs on a given index lets expose this to match the portal
			"default_ttl_seconds": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				ValidateFunc: validation.All(
					validation.IntAtLeast(-1),
					validation.IntNotInSlice([]int{0}),
				),
			},

			"analytical_storage_ttl": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(-1),
			},

			"throughput": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.CosmosThroughput,
			},

			"autoscale_settings": common.DatabaseAutoscaleSettingsSchema(),

			"index": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"keys": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},

						"unique": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"system_indexes": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"keys": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},

						"unique": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceCosmosDbMongoCollectionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := cosmosdb.NewMongodbDatabaseCollectionID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("database_name").(string), d.Get("name").(string))

	existing, err := client.MongoDBResourcesGetMongoDBCollection(ctx, id)
	if !response.WasNotFound(existing.HttpResponse) {
		if err != nil {
			return fmt.Errorf("checking for presence of %s: %+v", id, err)
		}
		return tf.ImportAsExistsError("azurerm_cosmosdb_mongo_collection", id.ID())
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

	if err := client.MongoDBResourcesCreateUpdateMongoDBCollectionThenPoll(ctx, id, db); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosDbMongoCollectionRead(d, meta)
}

func resourceCosmosDbMongoCollectionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
			return fmt.Errorf("index with '_id' key is required")
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

func resourceCosmosDbMongoCollectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
	if !isServerlessCapacityMode(accResp.Model) {
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

func resourceCosmosDbMongoCollectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseMongodbDatabaseCollectionID(d.Id())
	if err != nil {
		return err
	}

	if err := client.MongoDBResourcesDeleteMongoDBCollectionThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandCosmosMongoCollectionIndex(indexes []interface{}, defaultTtl *int) (*[]cosmosdb.MongoIndex, bool) {
	results := make([]cosmosdb.MongoIndex, 0)

	hasIdKey := false

	if len(indexes) != 0 {
		for _, v := range indexes {
			index := v.(map[string]interface{})
			keys := index["keys"].([]interface{})

			for _, key := range keys {
				if strings.EqualFold("_id", key.(string)) {
					hasIdKey = true
				}
			}

			results = append(results, cosmosdb.MongoIndex{
				Key: &cosmosdb.MongoIndexKeys{
					Keys: utils.ExpandStringSlice(index["keys"].([]interface{})),
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

func flattenCosmosMongoCollectionIndex(input *[]cosmosdb.MongoIndex, accountIsVersion36 bool) (*[]map[string]interface{}, *[]map[string]interface{}, *int64) {
	indexes := make([]map[string]interface{}, 0)
	systemIndexes := make([]map[string]interface{}, 0)
	var ttl *int64
	if input == nil {
		return &indexes, &systemIndexes, ttl
	}

	for _, v := range *input {
		index := map[string]interface{}{}
		systemIndex := map[string]interface{}{}

		if v.Key != nil && v.Key.Keys != nil && len(*v.Key.Keys) > 0 {
			key := (*v.Key.Keys)[0]

			switch key {
			// As `DocumentDBDefaultIndex` and `_id` cannot be updated, so they would be moved into `system_indexes`.
			case "_id":
				systemIndex["keys"] = utils.FlattenStringSlice(v.Key.Keys)
				// The system index `_id` is always unique but api returns nil and it would be converted to `false` by zero-value. So it has to be manually set as `true`.
				systemIndex["unique"] = true

				systemIndexes = append(systemIndexes, systemIndex)

				if accountIsVersion36 {
					index["keys"] = utils.FlattenStringSlice(v.Key.Keys)
					index["unique"] = true
					indexes = append(indexes, index)
				}

			case "DocumentDBDefaultIndex":
				// Updating system index `DocumentDBDefaultIndex` is not a supported scenario.
				systemIndex["keys"] = utils.FlattenStringSlice(v.Key.Keys)

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
				index["keys"] = utils.FlattenStringSlice(v.Key.Keys)

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
