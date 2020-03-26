package cosmos

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2019-08-01/documentdb"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/migration"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmCosmosDbMongoCollection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCosmosDbMongoCollectionCreate,
		Read:   resourceArmCosmosDbMongoCollectionRead,
		Update: resourceArmCosmosDbMongoCollectionUpdate,
		Delete: resourceArmCosmosDbMongoCollectionDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    migration.ResourceMongoDbCollectionUpgradeV0Schema().CoreConfigSchema().ImpliedType(),
				Upgrade: migration.ResourceMongoDbCollectionStateUpgradeV0ToV1,
				Version: 0,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosEntityName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosAccountName,
			},

			"database_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosEntityName,
			},

			// SDK/api accepts an array.. but only one is allowed
			"shard_key": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// default TTL is simply an index on _ts with expireAfterOption, given we can't seem to set TTLs on a given index lets expose this to match the portal
			"default_ttl_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(-1),
			},

			"throughput": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.CosmosThroughput,
			},
		},
	}
}

func resourceArmCosmosDbMongoCollectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.MongoDbClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	account := d.Get("account_name").(string)
	database := d.Get("database_name").(string)

	if features.ShouldResourcesBeImported() {
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
	}

	var ttl *int
	if v, ok := d.GetOkExists("default_ttl_seconds"); ok {
		ttl = utils.Int(v.(int))
	}

	db := documentdb.MongoDBCollectionCreateUpdateParameters{
		MongoDBCollectionCreateUpdateProperties: &documentdb.MongoDBCollectionCreateUpdateProperties{
			Resource: &documentdb.MongoDBCollectionResource{
				ID:      &name,
				Indexes: expandCosmosMongoCollectionIndexes(ttl),
			},
			Options: map[string]*string{},
		},
	}

	if throughput, hasThroughput := d.GetOk("throughput"); hasThroughput {
		db.MongoDBCollectionCreateUpdateProperties.Options = map[string]*string{
			"throughput": utils.String(strconv.Itoa(throughput.(int))),
		}
	}

	if v, ok := d.GetOkExists("shard_key"); ok {
		db.MongoDBCollectionCreateUpdateProperties.Resource.ShardKey = map[string]*string{
			v.(string): utils.String("Hash"), // looks like only hash is supported for now
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

	return resourceArmCosmosDbMongoCollectionRead(d, meta)
}

func resourceArmCosmosDbMongoCollectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.MongoDbClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MongoDbCollectionID(d.Id())
	if err != nil {
		return err
	}

	var ttl *int
	if v, ok := d.GetOkExists("default_ttl_seconds"); ok {
		ttl = utils.Int(v.(int))
	}

	db := documentdb.MongoDBCollectionCreateUpdateParameters{
		MongoDBCollectionCreateUpdateProperties: &documentdb.MongoDBCollectionCreateUpdateProperties{
			Resource: &documentdb.MongoDBCollectionResource{
				ID:      &id.Name,
				Indexes: expandCosmosMongoCollectionIndexes(ttl),
			},
			Options: map[string]*string{},
		},
	}

	if v, ok := d.GetOkExists("shard_key"); ok {
		db.MongoDBCollectionCreateUpdateProperties.Resource.ShardKey = map[string]*string{
			v.(string): utils.String("Hash"), // looks like only hash is supported for now
		}
	}

	future, err := client.CreateUpdateMongoDBCollection(ctx, id.ResourceGroup, id.Account, id.Database, id.Name, db)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Cosmos Mongo Collection %q (Account: %q, Database: %q): %+v", id.Name, id.Account, id.Database, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Cosmos Mongo Collection %q (Account: %q, Database: %q): %+v", id.Name, id.Account, id.Database, err)
	}

	if d.HasChange("throughput") {
		throughputParameters := documentdb.ThroughputSettingsUpdateParameters{
			ThroughputSettingsUpdateProperties: &documentdb.ThroughputSettingsUpdateProperties{
				Resource: &documentdb.ThroughputSettingsResource{
					Throughput: utils.Int32(int32(d.Get("throughput").(int))),
				},
			},
		}

		throughputFuture, err := client.UpdateMongoDBCollectionThroughput(ctx, id.ResourceGroup, id.Account, id.Database, id.Name, throughputParameters)
		if err != nil {
			if response.WasNotFound(throughputFuture.Response()) {
				return fmt.Errorf("Error setting Throughput for Cosmos MongoDB Collection %q (Account: %q, Database: %q): %+v - "+
					"If the collection has not been created with an initial throughput, you cannot configure it later.", id.Name, id.Account, id.Database, err)
			}
		}

		if err = throughputFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting on ThroughputUpdate future for Cosmos Mongo Collection %q (Account: %q, Database: %q): %+v", id.Name, id.Account, id.Database, err)
		}
	}

	return resourceArmCosmosDbMongoCollectionRead(d, meta)
}

func resourceArmCosmosDbMongoCollectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.MongoDbClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MongoDbCollectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetMongoDBCollection(ctx, id.ResourceGroup, id.Account, id.Database, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Cosmos Mongo Collection %q (Account: %q, Database: %q)", id.Name, id.Account, id.Database)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Cosmos Mongo Collection %q (Account: %q, Database: %q): %+v", id.Name, id.Account, id.Database, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_name", id.Account)
	d.Set("database_name", id.Database)
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

			if res.Indexes != nil {
				d.Set("default_ttl_seconds", flattenCosmosMongoCollectionIndexes(res.Indexes))
			}
		}
	}

	throughputResp, err := client.GetMongoDBCollectionThroughput(ctx, id.ResourceGroup, id.Account, id.Database, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(throughputResp.Response) {
			return fmt.Errorf("Error reading Throughput on Cosmos Mongo Collection %q (Account: %q, Database: %q): %+v", id.Name, id.Account, id.Database, err)
		} else {
			d.Set("throughput", nil)
		}
	} else {
		d.Set("throughput", common.GetThroughputFromResult(throughputResp))
	}

	return nil
}

func resourceArmCosmosDbMongoCollectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.MongoDbClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MongoDbCollectionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.DeleteMongoDBCollection(ctx, id.ResourceGroup, id.Account, id.Database, id.Name)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting Cosmos Mongo Collection %q (Account: %q, Database: %q): %+v", id.Name, id.Account, id.Database, err)
		}
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting on delete future for Cosmos Mongo Collection %q (Account: %q, Database: %q): %+v", id.Name, id.Account, id.Database, err)
	}

	return nil
}

func expandCosmosMongoCollectionIndexes(defaultTtl *int) *[]documentdb.MongoIndex {
	outputs := make([]documentdb.MongoIndex, 0)

	if defaultTtl != nil {
		outputs = append(outputs, documentdb.MongoIndex{
			Key: &documentdb.MongoIndexKeys{
				Keys: &[]string{"_ts"},
			},
			Options: &documentdb.MongoIndexOptions{
				ExpireAfterSeconds: utils.Int32(int32(*defaultTtl)),
			},
		})
	}

	return &outputs
}

func flattenCosmosMongoCollectionIndexes(indexes *[]documentdb.MongoIndex) *int {
	var ttl int
	for _, i := range *indexes {
		if key := i.Key; key != nil {
			var ttlInner int32

			if keys := key.Keys; keys != nil && len(*keys) > 0 {
				k := (*keys)[0]

				if k == "_ts" {
					ttl = int(ttlInner)
				}
			}
		}
	}

	return &ttl
}
