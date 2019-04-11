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

func resourceArmCosmosCassandraKeyspace() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCosmosCassandraKeyspaceCreate,
		Read:   resourceArmCosmosCassandraKeyspaceRead,
		Delete: resourceArmCosmosCassandraKeyspaceDelete,

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

func resourceArmCosmosCassandraKeyspaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmosAccountsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	account := d.Get("account_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.GetCassandraKeyspace(ctx, resourceGroup, account, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of creating Cosmos Cassandra Keyspace %s (Account %s): %+v", name, account, err)
			}
		} else {
			id, err := azure.CosmosGetIDFromResponse(existing.Response)
			if err != nil {
				return fmt.Errorf("Error checking for presence of Cosmos Cassandra Keyspace '%s' (Account %s): %v", name, account, err)
			}

			return tf.ImportAsExistsError("azurerm_cosmos_cassandra_keyspace", id)
		}
	}

	db := documentdb.CassandraKeyspaceCreateUpdateParameters{
		CassandraKeyspaceCreateUpdateProperties: &documentdb.CassandraKeyspaceCreateUpdateProperties{
			Resource: &documentdb.CassandraKeyspaceResource{
				ID: &name,
			},
			Options: map[string]*string{},
		},
	}

	future, err := client.CreateCassandraKeyspace(ctx, resourceGroup, account, db)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Cosmos Cassandra Keyspace %s (Account %s): %+v", name, account, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Cosmos Cassandra Keyspace %s (Account %s): %+v", name, account, err)
	}

	//so  well resp.ID is not set...   guess instead grab the ID from response.request.url.path
	resp, err := client.GetCassandraKeyspace(ctx, resourceGroup, account, name)
	if err != nil {
		return fmt.Errorf("Error making get request for Cosmos Cassandra Keyspace %s (Account %s): %+v", name, account, err)
	}

	id, err := azure.CosmosGetIDFromResponse(resp.Response)
	if err != nil {
		return fmt.Errorf("Error creating Cosmos Cassandra Keyspace '%s' (Account %s) ID", name, account)
	}
	d.SetId(id)

	return resourceArmCosmosCassandraKeyspaceRead(d, meta)
}

func resourceArmCosmosCassandraKeyspaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmosAccountsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseCosmosKeyspaceResourceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetCassandraKeyspace(ctx, id.ResourceGroup, id.Account, id.Keyspace)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Cosmos Cassandra Keyspace %s (Account %s) - removing from state", id.Keyspace, id.Account)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Cosmos Cassandra Keyspace %s (Account %s): %+v", id.Keyspace, id.Account, err)
	}

	if props := resp.CassandraKeyspaceProperties; props != nil {
		d.Set("name", props.ID)
		d.Set("resource_group_name", id.ResourceGroup)
		d.Set("account_name", id.Account)
	}

	return nil
}

func resourceArmCosmosCassandraKeyspaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmosAccountsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseCosmosKeyspaceResourceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.DeleteCassandraKeyspace(ctx, id.ResourceGroup, id.Account, id.Keyspace)
	if err != nil {
		if !response.WasNotFound(resp.Response()) {
			return fmt.Errorf("Error deleting Cosmos Cassandra Keyspace %s (Account %s): %+v", id.Keyspace, id.Account, err)
		}
	}

	return nil
}
