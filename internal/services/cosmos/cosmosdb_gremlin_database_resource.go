// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosGremlinDatabase() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCosmosGremlinDatabaseCreate,
		Update: resourceCosmosGremlinDatabaseUpdate,
		Read:   resourceCosmosGremlinDatabaseRead,
		Delete: resourceCosmosGremlinDatabaseDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := cosmosdb.ParseGremlinDatabaseID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.GremlinDatabaseV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosEntityName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosAccountName,
			},

			"throughput": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.CosmosThroughput,
			},

			"autoscale_settings": common.DatabaseAutoscaleSettingsSchema(),
		},
	}
}

func resourceCosmosGremlinDatabaseCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := cosmosdb.NewGremlinDatabaseID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("name").(string))

	existing, err := client.GremlinResourcesGetGremlinDatabase(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_cosmosdb_gremlin_database", id.ID())
	}

	db := cosmosdb.GremlinDatabaseCreateUpdateParameters{
		Properties: cosmosdb.GremlinDatabaseCreateUpdateProperties{
			Resource: cosmosdb.GremlinDatabaseResource{
				Id: id.GremlinDatabaseName,
			},
			Options: &cosmosdb.CreateUpdateOptions{},
		},
	}

	if throughput, hasThroughput := d.GetOk("throughput"); hasThroughput {
		if throughput != 0 {
			db.Properties.Options.Throughput = common.ConvertThroughputFromResourceData(throughput)
		}
	}

	if _, hasAutoscaleSettings := d.GetOk("autoscale_settings"); hasAutoscaleSettings {
		db.Properties.Options.AutoScaleSettings = common.ExpandCosmosDbAutoscaleSettings(d)
	}

	err = client.GremlinResourcesCreateUpdateGremlinDatabaseThenPoll(ctx, id, db)
	if err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosGremlinDatabaseRead(d, meta)
}

func resourceCosmosGremlinDatabaseUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseGremlinDatabaseID(d.Id())
	if err != nil {
		return err
	}

	err = common.CheckForChangeFromAutoscaleAndManualThroughput(d)
	if err != nil {
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

	err = client.GremlinResourcesCreateUpdateGremlinDatabaseThenPoll(ctx, *id, db)
	if err != nil {
		return fmt.Errorf("updating %q: %+v", id, err)
	}

	if common.HasThroughputChange(d) {
		if err := client.GremlinResourcesUpdateGremlinDatabaseThroughputThenPoll(ctx, *id, common.ExpandCosmosDBThroughputSettingsUpdateParameters(d)); err != nil {
			return fmt.Errorf("setting Throughput for %s: %+v - If the collection has not been created with and initial throughput, you cannot configure it later", id, err)
		}
	}

	return resourceCosmosGremlinDatabaseRead(d, meta)
}

func resourceCosmosGremlinDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseGremlinDatabaseID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GremlinResourcesGetGremlinDatabase(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("account_name", id.DatabaseAccountName)
	d.Set("name", id.GremlinDatabaseName)

	databaseAccountID := cosmosdb.NewDatabaseAccountID(id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName)
	accResp, err := client.DatabaseAccountsGet(ctx, databaseAccountID)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", databaseAccountID, err)
	}

	if !isServerlessCapacityMode(accResp.Model) {
		throughputResp, err := client.GremlinResourcesGetGremlinDatabaseThroughput(ctx, *id)
		if err != nil {
			if !response.WasNotFound(throughputResp.HttpResponse) {
				return fmt.Errorf("retrieving Throughput for %s: %+v", id, err)
			} else {
				d.Set("throughput", nil)
				d.Set("autoscale_settings", nil)
			}
		} else {
			common.SetResourceDataThroughputFromResponse(*throughputResp.Model, d)
		}
	}

	return nil
}

func resourceCosmosGremlinDatabaseDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseGremlinDatabaseID(d.Id())
	if err != nil {
		return err
	}

	err = client.GremlinResourcesDeleteGremlinDatabaseThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
