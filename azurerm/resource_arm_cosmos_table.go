package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmCosmosTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCosmosTableCreate,
		Read:   resourceArmCosmosTableRead,
		Delete: resourceArmCosmosTableDelete,

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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosAccountName,
			},
		},
	}
}

func resourceArmCosmosTableCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmosAccountsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	account := d.Get("account_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.GetTable(ctx, resourceGroup, account, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of creating Cosmos Table %s (Account %s): %+v", name, account, err)
			}
		} else {
			id, err := azure.CosmosGetIDFromResponse(existing.Response)
			if err != nil {
				return fmt.Errorf("Error checking for presence of Cosmos Table '%s' (Account %s): %v", name, account, err)
			}

			return tf.ImportAsExistsError("azurerm_cosmos_mongo_database", id)
		}
	}

	db := documentdb.TableCreateUpdateParameters{
		TableCreateUpdateProperties: &documentdb.TableCreateUpdateProperties{
			Resource: &documentdb.TableResource{
				ID: &name,
			},
			Options: map[string]*string{},
		},
	}

	future, err := client.CreateTable(ctx, resourceGroup, account, db)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Cosmos Table %s (Account %s): %+v", name, account, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Cosmos Table %s (Account %s): %+v", name, account, err)
	}

	//so  well resp.ID is not set...   guess instead grab the ID from response.request.url.path
	resp, err := client.GetTable(ctx, resourceGroup, account, name)
	if err != nil {
		return fmt.Errorf("Error making get request for Cosmos Table %s (Account %s): %+v", name, account, err)
	}

	id, err := azure.CosmosGetIDFromResponse(resp.Response)
	if err != nil {
		return fmt.Errorf("Error creating Cosmos Table '%s' (Account %s) ID", name, account)
	}
	d.SetId(id)

	return resourceArmCosmosTableRead(d, meta)
}

func resourceArmCosmosTableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmosAccountsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseCosmosTableResourceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetTable(ctx, id.ResourceGroup, id.Account, id.Table)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Cosmos Table %s (Account %s) - removing from state", id.Table, id.Account)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Cosmos Table %s (Account %s): %+v", id.Table, id.Account, err)
	}

	if props := resp.TableProperties; props != nil {
		d.Set("name", props.ID)
		d.Set("resource_group_name", id.ResourceGroup)
		d.Set("account_name", id.Account)
	}

	return nil
}

func resourceArmCosmosTableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmosAccountsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseCosmosTableResourceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.DeleteTable(ctx, id.ResourceGroup, id.Account, id.Table)
	if err != nil {
		if !response.WasNotFound(resp.Response()) {
			return fmt.Errorf("Error deleting Cosmos Table %s (Account %s): %+v", id.Table, id.Account, err)
		}
	}

	return nil
}
