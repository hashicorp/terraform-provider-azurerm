package cosmos

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-01-15/documentdb"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/migration"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceCosmosDbMongoCollection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCosmosDbMongoCollectionCreate,
		Read:   resourceCosmosDbMongoCollectionRead,
		Update: resourceCosmosDbMongoCollectionUpdate,
		Delete: resourceCosmosDbMongoCollectionDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

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

			"resource_group_name": azure.SchemaResourceGroupName(),

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
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(-1),
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

			"autoscale_settings": common.MongoCollectionAutoscaleSettingsSchema(),

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
	client := meta.(*clients.Client).Cosmos.MongoDbClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	account := d.Get("account_name").(string)
	database := d.Get("database_name").(string)

	existing, err := client.GetMongoDBCollection(ctx, resourceGroup, account, database, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of creating Cosmos Mongo Collection %q (Account: %q, Database: %q): %+v", name, account, database, err)
		}
	} else {
		if existing.ID == nil && *existing.ID == "" {
			return fmt.Errorf("Error generating import ID for Cosmos Mongo Collection %q (Account: %q, Database: %q)", name, account, database)
		}

		return tf.ImportAsExistsError("azurerm_cosmosdb_mongo_collection", *existing.ID)
	}

	var ttl *int
	if v := d.Get("default_ttl_seconds").(int); v > 0 {
		ttl = utils.Int(v)
	}

	db := documentdb.MongoDBCollectionCreateUpdateParameters{
		MongoDBCollectionCreateUpdateProperties: &documentdb.MongoDBCollectionCreateUpdateProperties{
			Resource: &documentdb.MongoDBCollectionResource{
				ID:      &name,
				Indexes: expandCosmosMongoCollectionIndex(d.Get("index").(*pluginsdk.Set).List(), ttl),
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	if analyticalStorageTTL, ok := d.GetOk("analytical_storage_ttl"); ok {
		db.MongoDBCollectionCreateUpdateProperties.Resource.AnalyticalStorageTTL = utils.Int32(int32(analyticalStorageTTL.(int)))
	}

	if throughput, hasThroughput := d.GetOk("throughput"); hasThroughput {
		if throughput != 0 {
			db.MongoDBCollectionCreateUpdateProperties.Options.Throughput = common.ConvertThroughputFromResourceData(throughput)
		}
	}

	if _, hasAutoscaleSettings := d.GetOk("autoscale_settings"); hasAutoscaleSettings {
		db.MongoDBCollectionCreateUpdateProperties.Options.AutoscaleSettings = common.ExpandCosmosDbAutoscaleSettings(d)
	}

	if shardKey := d.Get("shard_key").(string); shardKey != "" {
		db.MongoDBCollectionCreateUpdateProperties.Resource.ShardKey = map[string]*string{
			shardKey: utils.String("Hash"), // looks like only hash is supported for now
		}
	}

	future, err := client.CreateUpdateMongoDBCollection(ctx, resourceGroup, account, database, name, db)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Cosmos Mongo Collection %q (Account: %q, Database: %q): %+v", name, account, database, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Cosmos Mongo Collection %q (Account: %q, Database: %q): %+v", name, account, database, err)
	}

	resp, err := client.GetMongoDBCollection(ctx, resourceGroup, account, database, name)
	if err != nil {
		return fmt.Errorf("Error making get request for Cosmos Mongo Collection %q (Account: %q, Database: %q): %+v", name, account, database, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Error getting ID from Cosmos Mongo Collection %q (Account: %q, Database: %q)", name, account, database)
	}

	d.SetId(*resp.ID)

	return resourceCosmosDbMongoCollectionRead(d, meta)
}

func resourceCosmosDbMongoCollectionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.MongoDbClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MongodbCollectionID(d.Id())
	if err != nil {
		return err
	}

	err = common.CheckForChangeFromAutoscaleAndManualThroughput(d)
	if err != nil {
		return fmt.Errorf("Error updating Cosmos Mongo Collection %q (Account: %q, Database: %q): %+v", id.CollectionName, id.DatabaseAccountName, id.MongodbDatabaseName, err)
	}

	var ttl *int
	if v := d.Get("default_ttl_seconds").(int); v > 0 {
		ttl = utils.Int(v)
	}

	db := documentdb.MongoDBCollectionCreateUpdateParameters{
		MongoDBCollectionCreateUpdateProperties: &documentdb.MongoDBCollectionCreateUpdateProperties{
			Resource: &documentdb.MongoDBCollectionResource{
				ID:      &id.CollectionName,
				Indexes: expandCosmosMongoCollectionIndex(d.Get("index").(*pluginsdk.Set).List(), ttl),
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	if analyticalStorageTTL, ok := d.GetOk("analytical_storage_ttl"); ok {
		db.MongoDBCollectionCreateUpdateProperties.Resource.AnalyticalStorageTTL = utils.Int32(int32(analyticalStorageTTL.(int)))
	}

	if shardKey := d.Get("shard_key").(string); shardKey != "" {
		db.MongoDBCollectionCreateUpdateProperties.Resource.ShardKey = map[string]*string{
			shardKey: utils.String("Hash"), // looks like only hash is supported for now
		}
	}

	future, err := client.CreateUpdateMongoDBCollection(ctx, id.ResourceGroup, id.DatabaseAccountName, id.MongodbDatabaseName, id.CollectionName, db)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Cosmos Mongo Collection %q (Account: %q, Database: %q): %+v", id.CollectionName, id.DatabaseAccountName, id.MongodbDatabaseName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Cosmos Mongo Collection %q (Account: %q, Database: %q): %+v", id.CollectionName, id.DatabaseAccountName, id.MongodbDatabaseName, err)
	}

	if common.HasThroughputChange(d) {
		throughputParameters := common.ExpandCosmosDBThroughputSettingsUpdateParameters(d)
		throughputFuture, err := client.UpdateMongoDBCollectionThroughput(ctx, id.ResourceGroup, id.DatabaseAccountName, id.MongodbDatabaseName, id.CollectionName, *throughputParameters)
		if err != nil {
			if response.WasNotFound(throughputFuture.Response()) {
				return fmt.Errorf("Error setting Throughput for Cosmos MongoDB Collection %q (Account: %q, Database: %q): %+v - "+
					"If the collection has not been created with an initial throughput, you cannot configure it later.", id.CollectionName, id.DatabaseAccountName, id.MongodbDatabaseName, err)
			}
		}

		if err = throughputFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting on ThroughputUpdate future for Cosmos Mongo Collection %q (Account: %q, Database: %q): %+v", id.CollectionName, id.DatabaseAccountName, id.MongodbDatabaseName, err)
		}
	}

	return resourceCosmosDbMongoCollectionRead(d, meta)
}

func resourceCosmosDbMongoCollectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.MongoDbClient
	accClient := meta.(*clients.Client).Cosmos.DatabaseClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MongodbCollectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetMongoDBCollection(ctx, id.ResourceGroup, id.DatabaseAccountName, id.MongodbDatabaseName, id.CollectionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Cosmos Mongo Collection %q (Account: %q, Database: %q)", id.CollectionName, id.DatabaseAccountName, id.MongodbDatabaseName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Cosmos Mongo Collection %q (Account: %q, Database: %q): %+v", id.CollectionName, id.DatabaseAccountName, id.MongodbDatabaseName, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_name", id.DatabaseAccountName)
	d.Set("database_name", id.MongodbDatabaseName)

	accResp, err := accClient.Get(ctx, id.ResourceGroup, id.DatabaseAccountName)
	if err != nil {
		return fmt.Errorf("reading Cosmos Account %q : %+v", id.DatabaseAccountName, err)
	}
	if props := resp.MongoDBCollectionGetProperties; props != nil {
		if res := props.Resource; res != nil {
			d.Set("name", res.ID)

			// you can only have one
			if len(res.ShardKey) > 2 {
				return fmt.Errorf("unexpected number of shard keys: %d", len(res.ShardKey))
			}

			for k := range res.ShardKey {
				d.Set("shard_key", k)
			}
			accountIsVersion36 := false
			if accProps := accResp.DatabaseAccountGetProperties; accProps != nil {
				if capabilities := accProps.Capabilities; capabilities != nil {
					for _, v := range *capabilities {
						if v.Name != nil && *v.Name == "EnableMongo" {
							accountIsVersion36 = true
						}
					}
				}
			}

			indexes, systemIndexes, ttl := flattenCosmosMongoCollectionIndex(res.Indexes, accountIsVersion36)
			if err := d.Set("default_ttl_seconds", ttl); err != nil {
				return fmt.Errorf("failed to set `default_ttl_seconds`: %+v", err)
			}
			if err := d.Set("index", indexes); err != nil {
				return fmt.Errorf("failed to set `index`: %+v", err)
			}
			if err := d.Set("system_indexes", systemIndexes); err != nil {
				return fmt.Errorf("failed to set `system_indexes`: %+v", err)
			}

			d.Set("analytical_storage_ttl", res.AnalyticalStorageTTL)
		}
	}

	// if the cosmos account is serverless calling the get throughput api would yield an error
	if !isServerlessCapacityMode(accResp) {
		throughputResp, err := client.GetMongoDBCollectionThroughput(ctx, id.ResourceGroup, id.DatabaseAccountName, id.MongodbDatabaseName, id.CollectionName)
		if err != nil {
			if !utils.ResponseWasNotFound(throughputResp.Response) {
				return fmt.Errorf("Error reading Throughput on Cosmos Mongo Collection %q (Account: %q, Database: %q): %+v", id.CollectionName, id.DatabaseAccountName, id.MongodbDatabaseName, err)
			} else {
				d.Set("throughput", nil)
				d.Set("autoscale_settings", nil)
			}
		} else {
			common.SetResourceDataThroughputFromResponse(throughputResp, d)
		}
	}

	return nil
}

func resourceCosmosDbMongoCollectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.MongoDbClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MongodbCollectionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.DeleteMongoDBCollection(ctx, id.ResourceGroup, id.DatabaseAccountName, id.MongodbDatabaseName, id.CollectionName)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting Cosmos Mongo Collection %q (Account: %q, Database: %q): %+v", id.CollectionName, id.DatabaseAccountName, id.MongodbDatabaseName, err)
		}
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting on delete future for Cosmos Mongo Collection %q (Account: %q, Database: %q): %+v", id.CollectionName, id.DatabaseAccountName, id.MongodbDatabaseName, err)
	}

	return nil
}

func expandCosmosMongoCollectionIndex(indexes []interface{}, defaultTtl *int) *[]documentdb.MongoIndex {
	results := make([]documentdb.MongoIndex, 0)

	if len(indexes) != 0 {
		for _, v := range indexes {
			index := v.(map[string]interface{})

			results = append(results, documentdb.MongoIndex{
				Key: &documentdb.MongoIndexKeys{
					Keys: utils.ExpandStringSlice(index["keys"].([]interface{})),
				},
				Options: &documentdb.MongoIndexOptions{
					Unique: utils.Bool(index["unique"].(bool)),
				},
			})
		}
	}

	if defaultTtl != nil {
		results = append(results, documentdb.MongoIndex{
			Key: &documentdb.MongoIndexKeys{
				Keys: &[]string{"_ts"},
			},
			Options: &documentdb.MongoIndexOptions{
				ExpireAfterSeconds: utils.Int32(int32(*defaultTtl)),
			},
		})
	}

	return &results
}

func flattenCosmosMongoCollectionIndex(input *[]documentdb.MongoIndex, accountIsVersion36 bool) (*[]map[string]interface{}, *[]map[string]interface{}, *int32) {
	indexes := make([]map[string]interface{}, 0)
	systemIndexes := make([]map[string]interface{}, 0)
	var ttl *int32
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
