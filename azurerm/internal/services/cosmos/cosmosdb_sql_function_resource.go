package cosmos

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-01-15/documentdb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceCosmosDbSQLFunction() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmosDbSQLFunctionCreateUpdate,
		Read:   resourceCosmosDbSQLFunctionRead,
		Update: resourceCosmosDbSQLFunctionCreateUpdate,
		Delete: resourceCosmosDbSQLFunctionDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.SqlFunctionID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"container_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SqlContainerID,
			},

			"body": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}
func resourceCosmosDbSQLFunctionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Cosmos.SqlResourceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	containerId, _ := parse.SqlContainerID(d.Get("container_id").(string))
	body := d.Get("body").(string)

	id := parse.NewSqlFunctionID(subscriptionId, containerId.ResourceGroup, containerId.DatabaseAccountName, containerId.SqlDatabaseName, containerId.ContainerName, name)

	if d.IsNewResource() {
		existing, err := client.GetSQLUserDefinedFunction(ctx, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName, id.UserDefinedFunctionName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing CosmosDb SqlFunction %q: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_cosmosdb_sql_function", id.ID())
		}
	}

	createUpdateSqlUserDefinedFunctionParameters := documentdb.SQLUserDefinedFunctionCreateUpdateParameters{
		SQLUserDefinedFunctionCreateUpdateProperties: &documentdb.SQLUserDefinedFunctionCreateUpdateProperties{
			Resource: &documentdb.SQLUserDefinedFunctionResource{
				ID:   &name,
				Body: &body,
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}
	future, err := client.CreateUpdateSQLUserDefinedFunction(ctx, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName, id.UserDefinedFunctionName, createUpdateSqlUserDefinedFunctionParameters)
	if err != nil {
		return fmt.Errorf("creating/updating CosmosDb SqlFunction %q: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of the CosmosDb SqlFunction %q: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCosmosDbSQLFunctionRead(d, meta)
}

func resourceCosmosDbSQLFunctionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlFunctionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetSQLUserDefinedFunction(ctx, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName, id.UserDefinedFunctionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] CosmosDb SqlFunction %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving CosmosDb SqlFunction %q: %+v", id, err)
	}
	containerId := parse.NewSqlContainerID(id.SubscriptionId, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName)
	d.Set("name", id.UserDefinedFunctionName)
	d.Set("container_id", containerId.ID())
	if props := resp.SQLUserDefinedFunctionGetProperties; props != nil {
		if props.Resource != nil {
			d.Set("body", props.Resource.Body)
		}
	}
	return nil
}

func resourceCosmosDbSQLFunctionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlFunctionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.DeleteSQLUserDefinedFunction(ctx, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName, id.UserDefinedFunctionName)
	if err != nil {
		return fmt.Errorf("deleting CosmosDb SqlFunction %q: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of the CosmosDb SqlFunction %q: %+v", id, err)
	}
	return nil
}
