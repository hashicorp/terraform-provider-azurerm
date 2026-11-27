// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_cassandra_table

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbCassandraTableUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseCassandraKeyspaceTableID(d.Id())
	if err != nil {
		return err
	}

	if err := common.CheckForChangeFromAutoscaleAndManualThroughput(d); err != nil {
		return fmt.Errorf("checking `autoscale_settings` and `throughput` for %s: %w", id, err)
	}

	existing, err := client.CassandraResourcesGetCassandraTable(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %w", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: properties was nil", id)
	}

	if existing.Model.Properties.Resource == nil {
		return fmt.Errorf("retrieving %s: resource was nil", id)
	}

	table := cosmosdb.CassandraTableCreateUpdateParameters{
		Properties: cosmosdb.CassandraTableCreateUpdateProperties{
			Resource: cosmosdb.CassandraTableResource{
				Id:         id.TableName,
				Schema:     existing.Model.Properties.Resource.Schema,
				DefaultTtl: existing.Model.Properties.Resource.DefaultTtl,
			},
			Options: &cosmosdb.CreateUpdateOptions{},
		},
	}

	if d.HasChange("default_ttl") {
		table.Properties.Resource.DefaultTtl = pointer.To(int64(d.Get("default_ttl").(int)))

		if err := client.CassandraResourcesCreateUpdateCassandraTableThenPoll(ctx, *id, table); err != nil {
			return fmt.Errorf("updating %s: %+v", *id, err)
		}
	}

	if common.HasThroughputChange(d) {
		if err := client.CassandraResourcesUpdateCassandraTableThroughputThenPoll(ctx, *id, common.ExpandCosmosDBThroughputSettingsUpdateParameters(d)); err != nil {
			return fmt.Errorf("setting Throughput for %s: %+v - If the collection has not been created with an initial throughput, you cannot configure it later", id, err)
		}
	}

	return resourceCosmosDbCassandraTableRead(d, meta)
}
