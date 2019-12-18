package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
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
				ValidateFunc: validate.NoEmptyStrings,
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

			"indexes": {
				Type:       schema.TypeSet,
				Optional:   true,
				Deprecated: "Indexes are ignored unless they are the shared key so have been deprecated.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:         schema.TypeString, // this is a list in the SDK/API, however any more then a single value causes a 404
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"unique": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true, // portal defaults to true
						},

						// expire_after_seconds is only allowed on `_ts`:
						// Unable to parse request payload due to the following reason: 'The 'expireAfterSeconds' option is supported on '_ts' field only.
					},
				},
			},
		},
	}
}

func resourceArmCosmosDbMongoCollectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Cosmos.DatabaseClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	account := d.Get("account_name").(string)
	database := d.Get("database_name").(string)

	if features.ShouldResourcesBeImported() {
		existing, err := client.GetMongoDBCollection(ctx, resourceGroup, account, database, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of creating Cosmos Mongo Collection %s (Account %s, Database %s): %+v", name, account, database, err)
			}
		} else {
			id, err := azure.CosmosGetIDFromResponse(existing.Response)
			if err != nil {
				return fmt.Errorf("Error generating import ID for Cosmos Mongo Collection %s (Account %s, Database %s)", name, account, database)
			}

			return tf.ImportAsExistsError("azurerm_cosmosdb_mongo_collection", id)
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
				Indexes: expandCosmosMongoCollectionIndexes(d.Get("indexes"), ttl),
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
		return fmt.Errorf("Error issuing create/update request for Cosmos Mongo Collection %s (Account %s, Database %s): %+v", name, account, database, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Cosmos Mongo Collection %s (Account %s, Database %s): %+v", name, account, database, err)
	}

	resp, err := client.GetMongoDBCollection(ctx, resourceGroup, account, database, name)
	if err != nil {
		return fmt.Errorf("Error making get request for Cosmos Mongo Collection %s (Account %s, Database %s): %+v", name, account, database, err)
	}

	id, err := azure.CosmosGetIDFromResponse(resp.Response)
	if err != nil {
		return fmt.Errorf("Error getting ID for Cosmos Mongo Collection %s (Account %s, Database %s) ID: %v", name, account, database, err)
	}
	d.SetId(id)

	return resourceArmCosmosDbMongoCollectionRead(d, meta)
}

func resourceArmCosmosDbMongoCollectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Cosmos.DatabaseClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseCosmosDatabaseCollectionID(d.Id())
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
				ID:      &id.Collection,
				Indexes: expandCosmosMongoCollectionIndexes(d.Get("indexes"), ttl),
			},
			Options: map[string]*string{},
		},
	}

	if v, ok := d.GetOkExists("shard_key"); ok {
		db.MongoDBCollectionCreateUpdateProperties.Resource.ShardKey = map[string]*string{
			v.(string): utils.String("Hash"), // looks like only hash is supported for now
		}
	}

	future, err := client.CreateUpdateMongoDBCollection(ctx, id.ResourceGroup, id.Account, id.Database, id.Collection, db)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Cosmos Mongo Collection %s (Account %s, Database %s): %+v", id.Collection, id.Account, id.Database, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Cosmos Mongo Collection %s (Account %s, Database %s): %+v", id.Collection, id.Account, id.Database, err)
	}

	if d.HasChange("throughput") {
		throughputParameters := documentdb.ThroughputUpdateParameters{
			ThroughputUpdateProperties: &documentdb.ThroughputUpdateProperties{
				Resource: &documentdb.ThroughputResource{
					Throughput: utils.Int32(int32(d.Get("throughput").(int))),
				},
			},
		}

		throughputFuture, err := client.UpdateMongoDBCollectionThroughput(ctx, id.ResourceGroup, id.Account, id.Database, id.Collection, throughputParameters)
		if err != nil {
			if throughputFuture.Response().StatusCode == http.StatusNotFound {
				d.Set("throughput", nil)
				return fmt.Errorf("Error setting Throughput for Cosmos MongoDB Collection %s (Account %s, Database %s): %+v - "+
					"If the collection has not been created with an initial throughput, you cannot configure it later.", id.Collection, id.Account, id.Database, err)
			}
		}

		if err = throughputFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting on ThroughputUpdate future for Cosmos Mongo Collection %s (Account %s, Database %s): %+v", id.Collection, id.Account, id.Database, err)
		}
	}

	return resourceArmCosmosDbMongoCollectionRead(d, meta)
}

func resourceArmCosmosDbMongoCollectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Cosmos.DatabaseClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseCosmosDatabaseCollectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetMongoDBCollection(ctx, id.ResourceGroup, id.Account, id.Database, id.Collection)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Cosmos Mongo Collection %s (Account %s, Database %s)", id.Collection, id.Account, id.Database)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Cosmos Mongo Collection %s (Account %s, Database %s): %+v", id.Collection, id.Account, id.Database, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_name", id.Account)
	d.Set("database_name", id.Database)
	if props := resp.MongoDBCollectionProperties; props != nil {
		d.Set("name", props.ID)

		// you can only have one
		if len(props.ShardKey) > 2 {
			return fmt.Errorf("unexpected number of shard keys: %d", len(props.ShardKey))
		}

		for k := range props.ShardKey {
			d.Set("shard_key", k)
		}

		if props.Indexes != nil {
			indexes, ttl := flattenCosmosMongoCollectionIndexes(props.Indexes)
			d.Set("default_ttl_seconds", ttl)
			if err := d.Set("indexes", indexes); err != nil {
				return fmt.Errorf("Error setting `indexes`: %+v", err)
			}
		}
	}

	throughputResp, err := client.GetMongoDBCollectionThroughput(ctx, id.ResourceGroup, id.Account, id.Database, id.Collection)
	if err != nil {
		if !utils.ResponseWasNotFound(throughputResp.Response) {
			return fmt.Errorf("Error reading Throughput on Cosmos Mongo Collection %s (Account %s, Database %s): %+v", id.Collection, id.Account, id.Database, err)
		} else {
			d.Set("throughput", nil)
		}
	} else {
		d.Set("throughput", throughputResp.Throughput)
	}

	return nil
}

func resourceArmCosmosDbMongoCollectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Cosmos.DatabaseClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseCosmosDatabaseCollectionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.DeleteMongoDBCollection(ctx, id.ResourceGroup, id.Account, id.Database, id.Collection)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting Cosmos Mongo Collection %s (Account %s, Database %s): %+v", id.Collection, id.Account, id.Database, err)
		}
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting on delete future for Cosmos Mongo Collection %s (Account %s, Database %s): %+v", id.Collection, id.Account, id.Database, err)
	}

	return nil
}

func expandCosmosMongoCollectionIndexes(input interface{}, defaultTtl *int) *[]documentdb.MongoIndex {
	outputs := make([]documentdb.MongoIndex, 0)

	for _, i := range input.(*schema.Set).List() {
		b := i.(map[string]interface{})
		outputs = append(outputs, documentdb.MongoIndex{
			Key: &documentdb.MongoIndexKeys{
				Keys: &[]string{b["key"].(string)},
			},
			Options: &documentdb.MongoIndexOptions{
				Unique: utils.Bool(b["unique"].(bool)),
			},
		})
	}

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

func flattenCosmosMongoCollectionIndexes(indexes *[]documentdb.MongoIndex) (*[]map[string]interface{}, *int) {
	slice := make([]map[string]interface{}, 0)

	var ttl int
	for _, i := range *indexes {
		if key := i.Key; key != nil {
			m := map[string]interface{}{}
			var ttlInner int32

			if options := i.Options; options != nil {
				if v := options.Unique; v != nil {
					m["unique"] = *v
				} else {
					m["unique"] = false // todo required? API sends back nothing for false
				}

				if v := options.ExpireAfterSeconds; v != nil {
					ttlInner = *v
				}
			}

			if keys := key.Keys; keys != nil && len(*keys) > 0 {
				k := (*keys)[0]

				if !strings.HasPrefix(k, "_") && k != "DocumentDBDefaultIndex" { // lets ignore system properties?
					m["key"] = k

					// only append indexes with a non system key
					slice = append(slice, m)
				}

				if k == "_ts" {
					ttl = int(ttlInner)
				}
			}
		}
	}

	return &slice, &ttl
}
