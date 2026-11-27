// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_cassandra_table

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbCassandraTableRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseCassandraKeyspaceTableID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.CassandraResourcesGetCassandraTable(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("cassandra_keyspace_id", cosmosdb.NewCassandraKeyspaceID(id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.CassandraKeyspaceName).ID())
	d.Set("name", id.TableName)

	if respModel := resp.Model; respModel != nil {
		if props := respModel.Properties; props != nil {
			if res := props.Resource; res != nil {
				d.Set("default_ttl", res.DefaultTtl)
				d.Set("analytical_storage_ttl", pointer.From(res.AnalyticalStorageTtl))
				d.Set("schema", flattenTableSchema(res.Schema))
			}
		}
	}

	databaseAccountID := cosmosdb.NewDatabaseAccountID(id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName)
	accResp, err := client.DatabaseAccountsGet(ctx, databaseAccountID)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", databaseAccountID, err)
	}

	if !helpers.IsServerlessCapacityMode(accResp.Model) {
		throughputResp, err := client.CassandraResourcesGetCassandraTableThroughput(ctx, *id)
		if err != nil {
			if !response.WasNotFound(throughputResp.HttpResponse) {
				return fmt.Errorf("retrieving Throughput for %s: %+v", *id, err)
			} else {
				d.Set("throughput", nil)
				d.Set("autoscale_settings", nil)
			}
		} else {
			common.SetResourceDataThroughputFromResponse(pointer.From(throughputResp.Model), d)
		}
	}
	return nil
}
