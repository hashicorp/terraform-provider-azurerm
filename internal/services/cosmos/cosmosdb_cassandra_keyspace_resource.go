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

func resourceCosmosDbCassandraKeyspace() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCosmosDbCassandraKeyspaceCreate,
		Read:   resourceCosmosDbCassandraKeyspaceRead,
		Update: resourceCosmosDbCassandraKeyspaceUpdate,
		Delete: resourceCosmosDbCassandraKeyspaceDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.CassandraKeyspaceID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.CassandraKeyspaceV0ToV1{},
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

func resourceCosmosDbCassandraKeyspaceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewCassandraKeyspaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("name").(string))

	existing, err := client.GetCassandraKeyspace(ctx, id.ResourceGroup, id.DatabaseAccountName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of creating %s: %+v", id, err)
		}
	} else {
		if existing.ID == nil && *existing.ID == "" {
			return fmt.Errorf("generating import ID for %s", id)
		}

		return tf.ImportAsExistsError("azurerm_cosmosdb_cassandra_keyspace", *existing.ID)
	}

	db := documentdb.CassandraKeyspaceCreateUpdateParameters{
		CassandraKeyspaceCreateUpdateProperties: &documentdb.CassandraKeyspaceCreateUpdateProperties{
			Resource: &documentdb.CassandraKeyspaceResource{
				ID: &id.Name,
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	if throughput, hasThroughput := d.GetOk("throughput"); hasThroughput {
		if throughput != 0 {
			db.CassandraKeyspaceCreateUpdateProperties.Options.Throughput = common.ConvertThroughputFromResourceDataLegacy(throughput)
		}
	}

	if _, hasAutoscaleSettings := d.GetOk("autoscale_settings"); hasAutoscaleSettings {
		db.CassandraKeyspaceCreateUpdateProperties.Options.AutoscaleSettings = common.ExpandCosmosDbAutoscaleSettingsLegacy(d)
	}

	future, err := client.CreateUpdateCassandraKeyspace(ctx, id.ResourceGroup, id.DatabaseAccountName, id.Name, db)
	if err != nil {
		return fmt.Errorf("issuing create/update request for %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on create/update future for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosDbCassandraKeyspaceRead(d, meta)
}

func resourceCosmosDbCassandraKeyspaceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CassandraKeyspaceID(d.Id())
	if err != nil {
		return err
	}

	err = common.CheckForChangeFromAutoscaleAndManualThroughput(d)
	if err != nil {
		return fmt.Errorf("updating Cosmos Cassandra Keyspace %q (Account: %q) - %+v", id.Name, id.DatabaseAccountName, err)
	}

	db := documentdb.CassandraKeyspaceCreateUpdateParameters{
		CassandraKeyspaceCreateUpdateProperties: &documentdb.CassandraKeyspaceCreateUpdateProperties{
			Resource: &documentdb.CassandraKeyspaceResource{
				ID: &id.Name,
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	future, err := client.CreateUpdateCassandraKeyspace(ctx, id.ResourceGroup, id.DatabaseAccountName, id.Name, db)
	if err != nil {
		return fmt.Errorf("issuing create/update request for Cosmos Cassandra Keyspace %q (Account: %q): %+v", id.ResourceGroup, id.DatabaseAccountName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on create/update future for Cosmos Cassandra Keyspace %q (Account: %q): %+v", id.ResourceGroup, id.DatabaseAccountName, err)
	}

	if common.HasThroughputChange(d) {
		throughputParameters := common.ExpandCosmosDBThroughputSettingsUpdateParametersLegacy(d)
		throughputFuture, err := client.UpdateCassandraKeyspaceThroughput(ctx, id.ResourceGroup, id.DatabaseAccountName, id.Name, *throughputParameters)
		if err != nil {
			if response.WasNotFound(throughputFuture.Response()) {
				return fmt.Errorf("setting Throughput for Cosmos Cassandra Keyspace %q (Account: %q): %+v - "+
					"If the collection has not been created with an initial throughput, you cannot configure it later", id.Name, id.DatabaseAccountName, err)
			}
		}

		if err = throughputFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting on ThroughputUpdate future for Cosmos Cassandra Keyspace %q (Account: %q): %+v", id.Name, id.DatabaseAccountName, err)
		}
	}

	return resourceCosmosDbCassandraKeyspaceRead(d, meta)
}

func resourceCosmosDbCassandraKeyspaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClient
	accountClient := meta.(*clients.Client).Cosmos.DatabaseClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CassandraKeyspaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetCassandraKeyspace(ctx, id.ResourceGroup, id.DatabaseAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Cosmos Cassandra Keyspace %q (Account: %q) - removing from state", id.Name, id.DatabaseAccountName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Cosmos Cassandra Keyspace %q (Account: %q): %+v", id.Name, id.DatabaseAccountName, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_name", id.DatabaseAccountName)
	if props := resp.CassandraKeyspaceGetProperties; props != nil {
		if res := props.Resource; res != nil {
			d.Set("name", res.ID)
		}
	}

	accResp, err := accountClient.Get(ctx, id.ResourceGroup, id.DatabaseAccountName)
	if err != nil {
		return fmt.Errorf("reading Cosmos Account %q : %+v", id.DatabaseAccountName, err)
	}
	if accResp.ID == nil || *accResp.ID == "" {
		return fmt.Errorf("cosmosDB Account %q (Resource Group %q) ID is empty or nil", id.DatabaseAccountName, id.ResourceGroup)
	}

	if !isServerlessCapacityMode(accResp) {
		throughputResp, err := client.GetCassandraKeyspaceThroughput(ctx, id.ResourceGroup, id.DatabaseAccountName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(throughputResp.Response) {
				return fmt.Errorf("reading Throughput on Cosmos Cassandra Keyspace %q (Account: %q): %+v", id.Name, id.DatabaseAccountName, err)
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

func resourceCosmosDbCassandraKeyspaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CassandraKeyspaceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.DeleteCassandraKeyspace(ctx, id.ResourceGroup, id.DatabaseAccountName, id.Name)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("deleting Cosmos Cassandra Keyspace %q (Account: %q): %+v", id.Name, id.DatabaseAccountName, err)
		}
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("waiting on delete future for Cosmos Cassandra Keyspace %q (Account: %q): %+v", id.Name, id.DatabaseAccountName, err)
	}

	return nil
}
