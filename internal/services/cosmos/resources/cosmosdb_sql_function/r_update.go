// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_function

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbSQLFunctionUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseUserDefinedFunctionID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.SqlResourcesGetSqlUserDefinedFunction(ctx, *id); err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
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

	if err := client.SqlResourcesCreateUpdateSqlUserDefinedFunctionThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceCosmosDbSQLFunctionRead(d, meta)
}
