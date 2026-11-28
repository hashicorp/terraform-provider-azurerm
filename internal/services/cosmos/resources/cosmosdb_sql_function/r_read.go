// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_function

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbSQLFunctionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseUserDefinedFunctionID(d.Id())
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
	d.Set("container_id", cosmosdb.NewContainerID(id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName).ID())

	if resp.Model != nil {
		if props := resp.Model.Properties; props != nil && props.Resource != nil {
			d.Set("body", props.Resource.Body)
		}
	}

	return nil
}
