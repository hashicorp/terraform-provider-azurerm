// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_gremlin_database

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosGremlinDatabaseUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseGremlinDatabaseID(d.Id())
	if err != nil {
		return err
	}

	if err = common.CheckForChangeFromAutoscaleAndManualThroughput(d); err != nil {
		return fmt.Errorf("checking `autoscale_settings` and `throughput` for %s: %w", id, err)
	}

	db := cosmosdb.GremlinDatabaseCreateUpdateParameters{
		Properties: cosmosdb.GremlinDatabaseCreateUpdateProperties{
			Resource: cosmosdb.GremlinDatabaseResource{
				Id: id.GremlinDatabaseName,
			},
			Options: &cosmosdb.CreateUpdateOptions{},
		},
	}

	if err = client.GremlinResourcesCreateUpdateGremlinDatabaseThenPoll(ctx, *id, db); err != nil {
		return fmt.Errorf("updating %q: %+v", id, err)
	}

	if common.HasThroughputChange(d) {
		if err := client.GremlinResourcesUpdateGremlinDatabaseThroughputThenPoll(ctx, *id, common.ExpandCosmosDBThroughputSettingsUpdateParameters(d)); err != nil {
			return fmt.Errorf("setting Throughput for %s: %+v - If the collection has not been created with and initial throughput, you cannot configure it later", id, err)
		}
	}

	return resourceCosmosGremlinDatabaseRead(d, meta)
}
