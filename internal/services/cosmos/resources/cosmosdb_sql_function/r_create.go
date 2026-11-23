// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_function

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbSQLFunctionCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	containerId, err := cosmosdb.ParseContainerID(d.Get("container_id").(string))
	if err != nil {
		return err
	}

	id := cosmosdb.NewUserDefinedFunctionID(meta.(*clients.Client).Account.SubscriptionId, containerId.ResourceGroupName, containerId.DatabaseAccountName, containerId.SqlDatabaseName, containerId.ContainerName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.SqlResourcesGetSqlUserDefinedFunction(ctx, id)
		if !response.WasNotFound(existing.HttpResponse) {
			if err != nil {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
			return tf.ImportAsExistsError("azurerm_cosmosdb_sql_function", id.ID())
		}
	}

	payload := cosmosdb.SqlUserDefinedFunctionCreateUpdateParameters{
		Properties: cosmosdb.SqlUserDefinedFunctionCreateUpdateProperties{
			Resource: cosmosdb.SqlUserDefinedFunctionResource{
				Id:   id.UserDefinedFunctionName,
				Body: pointer.To(d.Get("body").(string)),
			},
			Options: &cosmosdb.CreateUpdateOptions{},
		},
	}

	if err := client.SqlResourcesCreateUpdateSqlUserDefinedFunctionCallbackThenPoll(ctx, id, payload, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCosmosDbSQLFunctionRead(d, meta)
}
