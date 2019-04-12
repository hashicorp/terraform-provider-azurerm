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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmCosmosCassandraTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCosmosCassandraTableCreateUpdate,
		Read:   resourceArmCosmosCassandraTableRead,
		Update: resourceArmCosmosCassandraTableCreateUpdate,
		Delete: resourceArmCosmosCassandraTableDelete,

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

			"resource_group_name": resourceGroupNameSchema(),

			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosAccountName,
			},

			"keyspace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosEntityName,
			},

			"default_ttl_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceArmCosmosCassandraTableCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmosAccountsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	account := d.Get("account_name").(string)
	keyspace := d.Get("keyspace_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.GetCassandraTable(ctx, resourceGroup, account, keyspace, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of creating Cosmos Cassandra Table %s (Account %s, Keyspace %s): %+v", name, account, keyspace, err)
			}
		} else {
			id, err := azure.CosmosGetIDFromResponse(existing.Response)
			if err != nil {
				return fmt.Errorf("Error generating import ID for Cosmos Cassandra Table %s (Account %s, Keyspace %s)", name, account, keyspace)
			}

			return tf.ImportAsExistsError("azurerm_cosmos_sql_container", id)
		}
	}

	db := documentdb.CassandraTableCreateUpdateParameters{
		CassandraTableCreateUpdateProperties: &documentdb.CassandraTableCreateUpdateProperties{
			Resource: &documentdb.CassandraTableResource{
				ID: &name,
			},
			Options: map[string]*string{},
		},
	}

	if d.IsNewResource() {
		future, err := client.CreateCassandraTable(ctx, resourceGroup, account, keyspace, db)
		if err != nil {
			return fmt.Errorf("Error issuing create request for Cosmos Cassandra Table %s (Account %s, Keyspace %s): %+v", name, account, keyspace, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting on create future for Cosmos Cassandra Table %s (Account %s, Keyspace %s): %+v", name, account, keyspace, err)
		}
	} else {
		future, err := client.UpdateCassandraTable(ctx, resourceGroup, account, keyspace, name, db)
		if err != nil {
			return fmt.Errorf("Error issuing update request for Cosmos Cassandra Table %s (Account %s, Keyspace %s): %+v", name, account, keyspace, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting on update future for Cosmos Cassandra Table %s (Account %s, Keyspace %s): %+v", name, account, keyspace, err)
		}
	}

	resp, err := client.GetCassandraTable(ctx, resourceGroup, account, keyspace, name)
	if err != nil {
		return fmt.Errorf("Error making get request for Cosmos Cassandra Table %s (Account %s, Keyspace %s): %+v", name, account, keyspace, err)
	}

	id, err := azure.CosmosGetIDFromResponse(resp.Response)
	if err != nil {
		return fmt.Errorf("Error creating Cosmos Cassandra Table %s (Account %s, Keyspace %s) ID: %v", name, account, keyspace, err)
	}
	d.SetId(id)

	return resourceArmCosmosCassandraTableRead(d, meta)
}

func resourceArmCosmosCassandraTableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmosAccountsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseCosmosKeyspaceTableID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetCassandraTable(ctx, id.ResourceGroup, id.Account, id.Keyspace, id.Table)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Cosmos Cassandra Table %s (Account %s, Keyspace %s)", id.Table, id.Account, id.Keyspace)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Cosmos Cassandra Table %s (Account %s, Keyspace %s): %+v", id.Table, id.Account, id.Keyspace, err)
	}

	if props := resp.CassandraTableProperties; props != nil {
		d.Set("name", props.ID)
		d.Set("resource_group_name", id.ResourceGroup)
		d.Set("account_name", id.Account)
		d.Set("keyspace_name", id.Keyspace)

	}

	return nil
}

func resourceArmCosmosCassandraTableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmosAccountsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseCosmosKeyspaceTableID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.DeleteCassandraTable(ctx, id.ResourceGroup, id.Account, id.Keyspace, id.Table)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting Cosmos Cassandra Table %s (Account %s, Keyspace %s): %+v", id.Table, id.Account, id.Keyspace, err)
		}
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting on delete future for Cosmos Cassandra Table %s (Account %s, Keyspace %s): %+v", id.Table, id.Account, id.Keyspace, err)
	}

	return nil
}
