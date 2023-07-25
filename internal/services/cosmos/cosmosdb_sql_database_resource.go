// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCosmosDbSQLDatabase() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCosmosDbSQLDatabaseCreate,
		Read:   resourceCosmosDbSQLDatabaseRead,
		Update: resourceCosmosDbSQLDatabaseUpdate,
		Delete: resourceCosmosDbSQLDatabaseDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SqlDatabaseID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SqlDatabaseV0ToV1{},
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

func resourceCosmosDbSQLDatabaseCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewSqlDatabaseID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("name").(string))

	existing, err := client.GetSQLDatabase(ctx, id.ResourceGroup, id.DatabaseAccountName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of %s: %+v", id, err)
		}
	} else {
		if existing.ID == nil && *existing.ID == "" {
			return fmt.Errorf("generating import ID for %s", id)
		}

		return tf.ImportAsExistsError("azurerm_cosmosdb_sql_database", *existing.ID)
	}

	db := documentdb.SQLDatabaseCreateUpdateParameters{
		SQLDatabaseCreateUpdateProperties: &documentdb.SQLDatabaseCreateUpdateProperties{
			Resource: &documentdb.SQLDatabaseResource{
				ID: &id.Name,
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	if throughput, hasThroughput := d.GetOk("throughput"); hasThroughput {
		if throughput != 0 {
			db.SQLDatabaseCreateUpdateProperties.Options.Throughput = common.ConvertThroughputFromResourceDataLegacy(throughput)
		}
	}

	if _, hasAutoscaleSettings := d.GetOk("autoscale_settings"); hasAutoscaleSettings {
		db.SQLDatabaseCreateUpdateProperties.Options.AutoscaleSettings = common.ExpandCosmosDbAutoscaleSettingsLegacy(d)
	}

	future, err := client.CreateUpdateSQLDatabase(ctx, id.ResourceGroup, id.DatabaseAccountName, id.Name, db)
	if err != nil {
		return fmt.Errorf("issuing create/update request for %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on create/update future for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosDbSQLDatabaseRead(d, meta)
}

func resourceCosmosDbSQLDatabaseUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlDatabaseID(d.Id())
	if err != nil {
		return err
	}

	err = common.CheckForChangeFromAutoscaleAndManualThroughput(d)
	if err != nil {
		return fmt.Errorf("updating Cosmos SQL Database %q (Account: %q) - %+v", id.Name, id.DatabaseAccountName, err)
	}

	db := documentdb.SQLDatabaseCreateUpdateParameters{
		SQLDatabaseCreateUpdateProperties: &documentdb.SQLDatabaseCreateUpdateProperties{
			Resource: &documentdb.SQLDatabaseResource{
				ID: &id.Name,
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	future, err := client.CreateUpdateSQLDatabase(ctx, id.ResourceGroup, id.DatabaseAccountName, id.Name, db)
	if err != nil {
		return fmt.Errorf("issuing create/update request for Cosmos SQL Database %q (Account: %q): %+v", id.Name, id.DatabaseAccountName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on create/update future for Cosmos SQL Database %q (Account: %q): %+v", id.Name, id.DatabaseAccountName, err)
	}

	if common.HasThroughputChange(d) {
		throughputParameters := common.ExpandCosmosDBThroughputSettingsUpdateParametersLegacy(d)
		throughputFuture, err := client.UpdateSQLDatabaseThroughput(ctx, id.ResourceGroup, id.DatabaseAccountName, id.Name, *throughputParameters)
		if err != nil {
			if response.WasNotFound(throughputFuture.Response()) {
				return fmt.Errorf("setting Throughput for Cosmos SQL Database %q (Account: %q) %+v - "+
					"If the collection has not been created with an initial throughput, you cannot configure it later.", id.Name, id.DatabaseAccountName, err)
			}
		}

		if err = throughputFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting on ThroughputUpdate future for Cosmos SQL Database %q (Account: %q): %+v", id.Name, id.DatabaseAccountName, err)
		}
	}

	return resourceCosmosDbSQLDatabaseRead(d, meta)
}

func resourceCosmosDbSQLDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlClient
	accountClient := meta.(*clients.Client).Cosmos.DatabaseClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlDatabaseID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetSQLDatabase(ctx, id.ResourceGroup, id.DatabaseAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Cosmos SQL Database %q (Account: %q) - removing from state", id.Name, id.DatabaseAccountName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Cosmos SQL Database %q (Account: %q): %+v", id.Name, id.DatabaseAccountName, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_name", id.DatabaseAccountName)
	if props := resp.SQLDatabaseGetProperties; props != nil {
		if res := props.Resource; res != nil {
			d.Set("name", res.ID)
		}
	}

	accResp, err := accountClient.Get(ctx, id.ResourceGroup, id.DatabaseAccountName)
	if err != nil {
		return fmt.Errorf("reading CosmosDB Account %q (Resource Group %q): %+v", id.DatabaseAccountName, id.ResourceGroup, err)
	}

	if accResp.ID == nil || *accResp.ID == "" {
		return fmt.Errorf("cosmosDB Account %q (Resource Group %q) ID is empty or nil", id.DatabaseAccountName, id.ResourceGroup)
	}

	// if the cosmos account is serverless calling the get throughput api would yield an error
	if !isServerlessCapacityMode(accResp) {
		throughputResp, err := client.GetSQLDatabaseThroughput(ctx, id.ResourceGroup, id.DatabaseAccountName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(throughputResp.Response) {
				return fmt.Errorf("reading Throughput on Cosmos SQL Database %q (Account: %q) ID: %v", id.Name, id.DatabaseAccountName, err)
			} else {
				d.Set("throughput", nil)
				d.Set("autoscale_settings", nil)
			}
		} else {
			common.SetResourceDataThroughputFromResponseLegacy(throughputResp, d)
		}
	}

	return nil
}

func resourceCosmosDbSQLDatabaseDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlDatabaseID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.DeleteSQLDatabase(ctx, id.ResourceGroup, id.DatabaseAccountName, id.Name)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("deleting Cosmos SQL Database %q (Account: %q): %+v", id.Name, id.DatabaseAccountName, err)
		}
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("waiting on delete future for Cosmos SQL Database %q (Account: %q): %+v", id.Name, id.DatabaseAccountName, err)
	}

	return nil
}
