// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_mongo_database

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

func resourceCosmosDbMongoDatabaseCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := cosmosdb.NewMongodbDatabaseID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.MongoDBResourcesGetMongoDBDatabase(ctx, id)
		if !response.WasNotFound(existing.HttpResponse) {
			if err != nil {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
			return tf.ImportAsExistsError("azurerm_cosmosdb_mongo_database", id.ID())
		}
	}

	db := cosmosdb.MongoDBDatabaseCreateUpdateParameters{
		Properties: cosmosdb.MongoDBDatabaseCreateUpdateProperties{
			Resource: cosmosdb.MongoDBDatabaseResource{
				Id: id.MongodbDatabaseName,
			},
			Options: &cosmosdb.CreateUpdateOptions{},
		},
	}

	if throughput, ok := d.GetOk("throughput"); ok {
		if throughput != 0 {
			db.Properties.Options.Throughput = common.ConvertThroughputFromResourceData(throughput)
		}
	}

	if _, ok := d.GetOk("autoscale_settings"); ok {
		db.Properties.Options.AutoScaleSettings = common.ExpandCosmosDbAutoscaleSettings(d)
	}

	if err := client.MongoDBResourcesCreateUpdateMongoDBDatabaseCallbackThenPoll(ctx, id, db, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosDbMongoDatabaseRead(d, meta)
}
