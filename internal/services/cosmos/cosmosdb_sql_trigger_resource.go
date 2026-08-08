// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2026-03-15/openapis"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
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
			_, err := openapis.ParseTriggerID(id)
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
				ValidateFunc: openapis.ValidateContainerID,
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
					string(openapis.TriggerOperationAll),
					string(openapis.TriggerOperationCreate),
					string(openapis.TriggerOperationUpdate),
					string(openapis.TriggerOperationDelete),
					string(openapis.TriggerOperationReplace),
				}, false),
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(openapis.TriggerTypePre),
					string(openapis.TriggerTypePost),
				}, false),
			},
		},
	}
}

func resourceCosmosDbSQLTriggerCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.OpenapisClient

	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	containerId, err := openapis.ParseContainerID(d.Get("container_id").(string))
	if err != nil {
		return err
	}

	id := openapis.NewTriggerID(meta.(*clients.Client).Account.SubscriptionId, containerId.ResourceGroupName, containerId.DatabaseAccountName, containerId.SqlDatabaseName, containerId.ContainerName, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.SqlResourcesGetSqlTrigger(ctx, id)
			if !response.WasNotFound(existing.HttpResponse) {
				if err != nil {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				return tf.ImportAsExistsError("azurerm_cosmosdb_sql_trigger", id.ID())
			}
		}
	}

	createUpdateSqlTriggerParameters := openapis.SqlTriggerCreateUpdateParameters{
		Properties: openapis.SqlTriggerCreateUpdateProperties{
			Resource: openapis.SqlTriggerResource{
				Id:               id.TriggerName,
				Body:             pointer.To(d.Get("body").(string)),
				TriggerType:      pointer.ToEnum[openapis.TriggerType](d.Get("type").(string)),
				TriggerOperation: pointer.ToEnum[openapis.TriggerOperation](d.Get("operation").(string)),
			},
			Options: &openapis.CreateUpdateOptions{},
		},
	}

	if d.IsNewResource() {
		if err := client.SqlResourcesCreateUpdateSqlTriggerCallbackThenPoll(ctx, id, createUpdateSqlTriggerParameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.SqlResourcesCreateUpdateSqlTriggerThenPoll(ctx, id, createUpdateSqlTriggerParameters); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceCosmosDbSQLTriggerRead(d, meta)
}

func resourceCosmosDbSQLTriggerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.OpenapisClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := openapis.ParseTriggerID(d.Id())
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
	d.Set("container_id", openapis.NewContainerID(id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName).ID())

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
	client := meta.(*clients.Client).Cosmos.OpenapisClient

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := openapis.ParseTriggerID(d.Id())
	if err != nil {
		return err
	}

	if err := client.SqlResourcesDeleteSqlTriggerThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
