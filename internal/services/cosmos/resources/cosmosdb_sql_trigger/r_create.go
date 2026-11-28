// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_trigger

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

func resourceCosmosDbSQLTriggerCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	containerId, err := cosmosdb.ParseContainerID(d.Get("container_id").(string))
	if err != nil {
		return err
	}

	id := cosmosdb.NewTriggerID(meta.(*clients.Client).Account.SubscriptionId, containerId.ResourceGroupName, containerId.DatabaseAccountName, containerId.SqlDatabaseName, containerId.ContainerName, d.Get("name").(string))

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
