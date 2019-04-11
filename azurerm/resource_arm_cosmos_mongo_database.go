package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmCosmosMongoDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCosmosMongoDatabaseCreate,
		Read:   resourceArmCosmosMongoDatabaseRead,
		Delete: resourceArmCosmosMongoDatabaseDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},

			"resource_group_name": resourceGroupNameSchema(),

			"account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-z0-9]{3,50}$"),
					"Cosmos DB Account name must be 3 - 50 characters long, contain only letters, numbers and hyphens.",
				),
			},
		},
	}
}

func resourceArmCosmosMongoDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmosAccountsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	account := d.Get("account_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.GetMongoDatabase(ctx, resourceGroup, account, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of creating Cosmos Mongo Database %s (Account %s): %+v", name, account, err)
			}
		} else {
			id, err := azure.CosmosGetIDFromResponse(existing.Response)
			if err != nil {
				return fmt.Errorf("Error checking for presence of CosmosDB Mongo Database '%s' (Account %s): %v", name, account, err)
			}

			return tf.ImportAsExistsError("azurerm_cosmos_mongo_database", id)
		}
	}

	db := documentdb.MongoDatabaseCreateUpdateParameters{
		MongoDatabaseCreateUpdateProperties: &documentdb.MongoDatabaseCreateUpdateProperties{
			Resource: &documentdb.MongoDatabaseResource{
				ID: &name,
			},
			Options: map[string]*string{},
		},
	}

	future, err := client.CreateMongoDatabase(ctx, resourceGroup, account, db)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Cosmos Mongo Database %s (Account %s): %+v", name, account, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Cosmos Mongo Database %s (Account %s): %+v", name, account, err)
	}

	//so  well resp.ID is not set...   grab the ID from response.request.url.path
	resp, err := client.GetMongoDatabase(ctx, resourceGroup, account, name)
	if err != nil {
		return fmt.Errorf("Error making get request for Cosmos Mongo Database %s (Account %s): %+v", name, account, err)
	}

	id, err := azure.CosmosGetIDFromResponse(resp.Response)
	if err != nil {
		return fmt.Errorf("Error creating CosmosDB Mongo Database '%s' (Account %s) ID", name, account)
	}
	d.SetId(id)

	return resourceArmCosmosMongoDatabaseRead(d, meta)
}

func resourceArmCosmosMongoDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmosAccountsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseCosmosDatabaseResourceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetMongoDatabase(ctx, id.ResourceGroup, id.Account, id.Database)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Cosmos Mongo Database %s (Account %s) - removing from state", id.Database, id.Account)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Cosmos Mongo Database %s (Account %s): %+v", id.Database, id.Account, err)
	}

	if props := resp.MongoDatabaseProperties; props != nil {
		d.Set("name", props.ID)
		d.Set("resource_group_name", id.ResourceGroup)
		d.Set("account_name", id.Account)
	}

	return nil
}

func resourceArmCosmosMongoDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmosAccountsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseCosmosDatabaseResourceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.DeleteMongoDatabase(ctx, id.ResourceGroup, id.Account, id.Database)
	if err != nil {
		if !response.WasNotFound(resp.Response()) {
			return fmt.Errorf("Error deleting Cosmos Mongo Database %s (Account %s): %+v", id.Database, id.Account, err)
		}
	}

	return nil
}
