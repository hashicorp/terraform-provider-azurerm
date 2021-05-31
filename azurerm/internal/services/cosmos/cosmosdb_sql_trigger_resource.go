package cosmos

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-01-15/documentdb"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceCosmosDbSQLTrigger() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCosmosDbSQLTriggerCreateUpdate,
		Read:   resourceCosmosDbSQLTriggerRead,
		Update: resourceCosmosDbSQLTriggerCreateUpdate,
		Delete: resourceCosmosDbSQLTriggerDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SqlTriggerID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosEntityName,
			},

			"container_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SqlContainerID,
			},

			"body": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"operation": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(documentdb.All),
					string(documentdb.Create),
					string(documentdb.Update),
					string(documentdb.Delete),
					string(documentdb.Replace),
				}, false),
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(documentdb.Pre),
					string(documentdb.Post),
				}, false),
			},
		},
	}
}
func resourceCosmosDbSQLTriggerCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Cosmos.SqlResourceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	containerId, _ := parse.SqlContainerID(d.Get("container_id").(string))
	body := d.Get("body").(string)
	triggerOperation := d.Get("operation").(string)
	triggerType := d.Get("type").(string)

	id := parse.NewSqlTriggerID(subscriptionId, containerId.ResourceGroup, containerId.DatabaseAccountName, containerId.SqlDatabaseName, containerId.ContainerName, name)

	if d.IsNewResource() {
		existing, err := client.GetSQLTrigger(ctx, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName, id.TriggerName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing CosmosDb SQLTrigger %q: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_cosmosdb_sql_trigger", id.ID())
		}
	}

	createUpdateSqlTriggerParameters := documentdb.SQLTriggerCreateUpdateParameters{
		SQLTriggerCreateUpdateProperties: &documentdb.SQLTriggerCreateUpdateProperties{
			Resource: &documentdb.SQLTriggerResource{
				ID:               &name,
				Body:             &body,
				TriggerType:      documentdb.TriggerType(triggerType),
				TriggerOperation: documentdb.TriggerOperation(triggerOperation),
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}
	future, err := client.CreateUpdateSQLTrigger(ctx, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName, name, createUpdateSqlTriggerParameters)
	if err != nil {
		return fmt.Errorf("creating/updating CosmosDb SQLTrigger %q: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of the CosmosDb SQLTrigger %q: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCosmosDbSQLTriggerRead(d, meta)
}

func resourceCosmosDbSQLTriggerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlTriggerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetSQLTrigger(ctx, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName, id.TriggerName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] CosmosDb SQLTrigger %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving CosmosDb SQLTrigger %q: %+v", id, err)
	}
	containerId := parse.NewSqlContainerID(id.SubscriptionId, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName)
	d.Set("name", id.TriggerName)
	d.Set("container_id", containerId.ID())
	if props := resp.SQLTriggerGetProperties; props != nil {
		if props.Resource != nil {
			d.Set("body", props.Resource.Body)
			d.Set("operation", props.Resource.TriggerOperation)
			d.Set("type", props.Resource.TriggerType)
		}
	}
	return nil
}

func resourceCosmosDbSQLTriggerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlTriggerID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.DeleteSQLTrigger(ctx, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName, id.TriggerName)
	if err != nil {
		return fmt.Errorf("deleting CosmosDb SQLResourcesSQLTrigger %q: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of the CosmosDb SQLResourcesSQLTrigger %q: %+v", id, err)
	}
	return nil
}
