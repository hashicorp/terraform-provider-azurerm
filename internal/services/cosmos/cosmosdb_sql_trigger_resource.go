// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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
			_, err := cosmosdb.ParseTriggerID(id)
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
				ValidateFunc: cosmosdb.ValidateContainerID,
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
					string(cosmosdb.TriggerOperationAll),
					string(cosmosdb.TriggerOperationCreate),
					string(cosmosdb.TriggerOperationUpdate),
					string(cosmosdb.TriggerOperationDelete),
					string(cosmosdb.TriggerOperationReplace),
				}, false),
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cosmosdb.TriggerTypePre),
					string(cosmosdb.TriggerTypePost),
				}, false),
			},
		},
	}
}

func resourceCosmosDbSQLTriggerCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	containerId, err := cosmosdb.ParseContainerID(d.Get("container_id").(string))
	if err != nil {
		return err
	}

	id := cosmosdb.NewTriggerID(meta.(*clients.Client).Account.SubscriptionId, containerId.ResourceGroupName, containerId.DatabaseAccountName, containerId.SqlDatabaseName, containerId.ContainerName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.SqlResourcesGetSqlTrigger(ctx, id)
		if !response.WasNotFound(existing.HttpResponse) {
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			return tf.ImportAsExistsError("azurerm_cosmosdb_sql_trigger", id.ID())
		}
	}

	createUpdateSqlTriggerParameters := cosmosdb.SqlTriggerCreateUpdateParameters{
		Properties: cosmosdb.SqlTriggerCreateUpdateProperties{
			Resource: cosmosdb.SqlTriggerResource{
				Id:               id.TriggerName,
				Body:             pointer.To(d.Get("body").(string)),
				TriggerType:      pointer.ToEnum[cosmosdb.TriggerType](d.Get("type").(string)),
				TriggerOperation: pointer.ToEnum[cosmosdb.TriggerOperation](d.Get("operation").(string)),
			},
			Options: &cosmosdb.CreateUpdateOptions{},
		},
	}

	if err := client.SqlResourcesCreateUpdateSqlTriggerThenPoll(ctx, id, createUpdateSqlTriggerParameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCosmosDbSQLTriggerRead(d, meta)
}

func resourceCosmosDbSQLTriggerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseTriggerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.SqlResourcesGetSqlTrigger(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.TriggerName)
	d.Set("container_id", cosmosdb.NewContainerID(id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName).ID())

	if resp.Model != nil {
		if props := resp.Model.Properties; props != nil {
			if r := props.Resource; r != nil {
				d.Set("body", r.Body)
				d.Set("operation", pointer.FromEnum(r.TriggerOperation))
				d.Set("type", pointer.FromEnum(r.TriggerType))
			}
		}
	}

	return nil
}

func resourceCosmosDbSQLTriggerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseTriggerID(d.Id())
	if err != nil {
		return err
	}

	if err := client.SqlResourcesDeleteSqlTriggerThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
