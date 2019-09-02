package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmCosmosDbMongoDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCosmosDbMongoDatabaseCreate,
		Read:   resourceArmCosmosDbMongoDatabaseRead,
		Delete: resourceArmCosmosDbMongoDatabaseDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
		},
	}
}

func resourceArmCosmosDbMongoDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmos.DatabaseClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	account := d.Get("account_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.GetMongoDBDatabase(ctx, resourceGroup, account, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of creating Cosmos Mongo Database %s (Account %s): %+v", name, account, err)
			}
		} else {
			id, err := azure.CosmosGetIDFromResponse(existing.Response)
			if err != nil {
				return fmt.Errorf("Error generating import ID for Cosmos Mongo Database '%s' (Account %s)", name, account)
			}

			return tf.ImportAsExistsError("azurerm_cosmosdb_mongo_database", id)
		}
	}

	db := documentdb.MongoDBDatabaseCreateUpdateParameters{
		MongoDBDatabaseCreateUpdateProperties: &documentdb.MongoDBDatabaseCreateUpdateProperties{
			Resource: &documentdb.MongoDBDatabaseResource{
				ID: &name,
			},
			Options: map[string]*string{},
		},
	}

	future, err := client.CreateUpdateMongoDBDatabase(ctx, resourceGroup, account, name, db)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Cosmos Mongo Database %s (Account %s): %+v", name, account, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Cosmos Mongo Database %s (Account %s): %+v", name, account, err)
	}

	resp, err := client.GetMongoDBDatabase(ctx, resourceGroup, account, name)
	if err != nil {
		return fmt.Errorf("Error making get request for Cosmos Mongo Database %s (Account %s): %+v", name, account, err)
	}

	id, err := azure.CosmosGetIDFromResponse(resp.Response)
	if err != nil {
		return fmt.Errorf("Error retrieving the ID for Cosmos Mongo Database '%s' (Account %s) ID: %v", name, account, err)
	}
	d.SetId(id)

	return resourceArmCosmosDbMongoDatabaseRead(d, meta)
}

func resourceArmCosmosDbMongoDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmos.DatabaseClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseCosmosDatabaseID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetMongoDBDatabase(ctx, id.ResourceGroup, id.Account, id.Database)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Cosmos Mongo Database %s (Account %s) - removing from state", id.Database, id.Account)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Cosmos Mongo Database %s (Account %s): %+v", id.Database, id.Account, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_name", id.Account)
	if props := resp.MongoDBDatabaseProperties; props != nil {
		d.Set("name", props.ID)
	}

	return nil
}

func resourceArmCosmosDbMongoDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmos.DatabaseClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseCosmosDatabaseID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.DeleteMongoDBDatabase(ctx, id.ResourceGroup, id.Account, id.Database)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting Cosmos Mongo Database %s (Account %s): %+v", id.Database, id.Account, err)
		}
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting on delete future for Cosmos Mongo Database %s (Account %s): %+v", id.Database, id.Account, err)
	}

	return nil
}
