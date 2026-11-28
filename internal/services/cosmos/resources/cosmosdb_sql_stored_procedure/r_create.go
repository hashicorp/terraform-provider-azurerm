// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_stored_procedure

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

func resourceCosmosDbSQLStoredProcedureCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := cosmosdb.NewStoredProcedureID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("database_name").(string), d.Get("container_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.SqlResourcesGetSqlStoredProcedure(ctx, id)
		if !response.WasNotFound(existing.HttpResponse) {
			if err != nil {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
			return tf.ImportAsExistsError("azurerm_cosmosdb_sql_stored_procedure", id.ID())
		}
	}

	storedProcParams := cosmosdb.SqlStoredProcedureCreateUpdateParameters{
		Properties: cosmosdb.SqlStoredProcedureCreateUpdateProperties{
			Resource: cosmosdb.SqlStoredProcedureResource{
				Id:   id.StoredProcedureName,
				Body: pointer.To(d.Get("body").(string)),
			},
			Options: &cosmosdb.CreateUpdateOptions{},
		},
	}

	if err := client.SqlResourcesCreateUpdateSqlStoredProcedureCallbackThenPoll(ctx, id, storedProcParams, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosDbSQLStoredProcedureRead(d, meta)
}
