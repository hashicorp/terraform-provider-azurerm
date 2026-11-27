// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_database

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbSQLDatabaseCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := cosmosdb.NewSqlDatabaseID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.SqlResourcesGetSqlDatabase(ctx, id)
		if !response.WasNotFound(existing.HttpResponse) {
			if err != nil {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
			return tf.ImportAsExistsError("azurerm_cosmosdb_sql_database", id.ID())
		}
	}

	db := cosmosdb.SqlDatabaseCreateUpdateParameters{
		Properties: cosmosdb.SqlDatabaseCreateUpdateProperties{
			Resource: cosmosdb.SqlDatabaseResource{
				Id: id.SqlDatabaseName,
			},
			Options: &cosmosdb.CreateUpdateOptions{},
		},
	}

	if throughput, ok := d.GetOk("throughput"); ok && throughput != 0 {
		db.Properties.Options.Throughput = common.ConvertThroughputFromResourceData(throughput)
	}

	if _, ok := d.GetOk("autoscale_settings"); ok {
		db.Properties.Options.AutoScaleSettings = common.ExpandCosmosDbAutoscaleSettings(d)
	}

	if err := client.SqlResourcesCreateUpdateSqlDatabaseCallbackThenPoll(ctx, id, db, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosDbSQLDatabaseRead(d, meta)
}
