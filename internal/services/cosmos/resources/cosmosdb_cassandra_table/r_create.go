// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_cassandra_table

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbCassandraTableCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	keyspaceId, err := cosmosdb.ParseCassandraKeyspaceID(d.Get("cassandra_keyspace_id").(string))
	if err != nil {
		return err
	}

	id := cosmosdb.NewCassandraKeyspaceTableID(meta.(*clients.Client).Account.SubscriptionId, keyspaceId.ResourceGroupName, keyspaceId.DatabaseAccountName, keyspaceId.CassandraKeyspaceName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.CassandraResourcesGetCassandraTable(ctx, id)
		if !response.WasNotFound(existing.HttpResponse) {
			if err != nil {
				return fmt.Errorf("checking for presence of existing %+v: %+v", id, err)
			}
			return tf.ImportAsExistsError("azurerm_cosmosdb_cassandra_table", id.ID())
		}
	}

	table := cosmosdb.CassandraTableCreateUpdateParameters{
		Properties: cosmosdb.CassandraTableCreateUpdateProperties{
			Options: &cosmosdb.CreateUpdateOptions{},
			Resource: cosmosdb.CassandraTableResource{
				Id:     id.TableName,
				Schema: expandTableSchema(d),
			},
		},
	}

	if defaultTTL, ok := d.GetOk("default_ttl"); ok {
		table.Properties.Resource.DefaultTtl = pointer.To(int64(defaultTTL.(int)))
	}

	if analyticalTTL, ok := d.GetOk("analytical_storage_ttl"); ok {
		table.Properties.Resource.AnalyticalStorageTtl = pointer.To(int64(analyticalTTL.(int)))
	}

	if throughput, hasThroughput := d.GetOk("throughput"); hasThroughput {
		if throughput != 0 {
			table.Properties.Options.Throughput = common.ConvertThroughputFromResourceData(throughput)
		}
	}

	if _, hasAutoscaleSettings := d.GetOk("autoscale_settings"); hasAutoscaleSettings {
		table.Properties.Options.AutoScaleSettings = common.ExpandCosmosDbAutoscaleSettings(d)
	}

	if err := client.CassandraResourcesCreateUpdateCassandraTableCallbackThenPoll(ctx, id, table, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosDbCassandraTableRead(d, meta)
}
