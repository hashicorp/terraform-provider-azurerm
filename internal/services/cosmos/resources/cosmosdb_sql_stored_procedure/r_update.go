// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_stored_procedure

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbSQLStoredProcedureUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseStoredProcedureID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.SqlResourcesGetSqlStoredProcedure(ctx, *id); err != nil {
		return fmt.Errorf("retrieving %s: %w", id, err)
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

	if err := client.SqlResourcesCreateUpdateSqlStoredProcedureThenPoll(ctx, *id, storedProcParams); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceCosmosDbSQLStoredProcedureRead(d, meta)
}
