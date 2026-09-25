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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbSQLFunction() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCosmosDbSQLFunctionCreate,
		Read:   resourceCosmosDbSQLFunctionRead,
		Update: resourceCosmosDbSQLFunctionUpdate,
		Delete: resourceCosmosDbSQLFunctionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := openapis.ParseUserDefinedFunctionID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
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
		},
	}
}

func resourceCosmosDbSQLFunctionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.OpenapisClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	containerId, err := openapis.ParseContainerID(d.Get("container_id").(string))
	if err != nil {
		return err
	}

	id := openapis.NewUserDefinedFunctionID(meta.(*clients.Client).Account.SubscriptionId, containerId.ResourceGroupName, containerId.DatabaseAccountName, containerId.SqlDatabaseName, containerId.ContainerName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.SqlResourcesGetSqlUserDefinedFunction(ctx, id)
		if !response.WasNotFound(existing.HttpResponse) {
			if err != nil {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
			return tf.ImportAsExistsError("azurerm_cosmosdb_sql_function", id.ID())
		}
	}

	payload := openapis.SqlUserDefinedFunctionCreateUpdateParameters{
		Properties: openapis.SqlUserDefinedFunctionCreateUpdateProperties{
			Resource: openapis.SqlUserDefinedFunctionResource{
				Id:   id.UserDefinedFunctionName,
				Body: pointer.To(d.Get("body").(string)),
			},
			Options: &openapis.CreateUpdateOptions{},
		},
	}

	if err := client.SqlResourcesCreateUpdateSqlUserDefinedFunctionCallbackThenPoll(ctx, id, payload, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCosmosDbSQLFunctionRead(d, meta)
}

func resourceCosmosDbSQLFunctionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.OpenapisClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := openapis.ParseUserDefinedFunctionID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.SqlResourcesGetSqlUserDefinedFunction(ctx, *id); err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	payload := openapis.SqlUserDefinedFunctionCreateUpdateParameters{
		Properties: openapis.SqlUserDefinedFunctionCreateUpdateProperties{
			Resource: openapis.SqlUserDefinedFunctionResource{
				Id:   id.UserDefinedFunctionName,
				Body: pointer.To(d.Get("body").(string)),
			},
			Options: &openapis.CreateUpdateOptions{},
		},
	}

	if err := client.SqlResourcesCreateUpdateSqlUserDefinedFunctionThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceCosmosDbSQLFunctionRead(d, meta)
}

func resourceCosmosDbSQLFunctionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.OpenapisClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := openapis.ParseUserDefinedFunctionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.SqlResourcesGetSqlUserDefinedFunction(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.UserDefinedFunctionName)
	d.Set("container_id", openapis.NewContainerID(id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName).ID())

	if resp.Model != nil {
		if props := resp.Model.Properties; props != nil && props.Resource != nil {
			d.Set("body", props.Resource.Body)
		}
	}

	return nil
}

func resourceCosmosDbSQLFunctionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.OpenapisClient

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := openapis.ParseUserDefinedFunctionID(d.Id())
	if err != nil {
		return err
	}

	if err := client.SqlResourcesDeleteSqlUserDefinedFunctionThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
