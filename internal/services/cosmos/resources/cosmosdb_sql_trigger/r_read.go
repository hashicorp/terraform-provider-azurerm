// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_trigger

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbSQLTriggerRead(d *pluginsdk.ResourceData, meta any) error {
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
