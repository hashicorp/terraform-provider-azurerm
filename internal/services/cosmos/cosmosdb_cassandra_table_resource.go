// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCosmosDbCassandraTable() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCosmosDbCassandraTableCreate,
		Read:   resourceCosmosDbCassandraTableRead,
		Update: resourceCosmosDbCassandraTableUpdate,
		Delete: resourceCosmosDbCassandraTableDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.CassandraTableID(id)
			return err
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

			"cassandra_keyspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CassandraKeyspaceID,
			},

			"default_ttl": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(-1),
			},

			"analytical_storage_ttl": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.IntBetween(-1, 2147483647),
					validation.IntNotInSlice([]int{0}),
				),
			},

			"schema": common.CassandraTableSchemaPropertySchema(),

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

func resourceCosmosDbCassandraTableCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyspaceId, err := parse.CassandraKeyspaceID(d.Get("cassandra_keyspace_id").(string))
	if err != nil {
		return fmt.Errorf("parsing Cassandra Keyspace ID: %+v", err)
	}
	account := keyspaceId.DatabaseAccountName
	keyspace := keyspaceId.Name
	resourceGroup := keyspaceId.ResourceGroup

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := parse.NewCassandraTableID(subscriptionId, resourceGroup, account, keyspace, name)
	existing, err := client.GetCassandraTable(ctx, resourceGroup, account, keyspace, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %+v: %+v", id, err)
		}
	} else {
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_cosmosdb_cassandra_table", id.ID())
		}
	}

	table := documentdb.CassandraTableCreateUpdateParameters{
		CassandraTableCreateUpdateProperties: &documentdb.CassandraTableCreateUpdateProperties{
			Resource: &documentdb.CassandraTableResource{
				ID: &name,
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	table.CassandraTableCreateUpdateProperties.Resource.Schema = expandTableSchema(d)

	if defaultTTL, hasTTL := d.GetOk("default_ttl"); hasTTL {
		table.CassandraTableCreateUpdateProperties.Resource.DefaultTTL = utils.Int32(int32(defaultTTL.(int)))
	}

	if analyticalTTL, ok := d.GetOk("analytical_storage_ttl"); ok {
		table.CassandraTableCreateUpdateProperties.Resource.AnalyticalStorageTTL = utils.Int32(int32(analyticalTTL.(int)))
	}

	if throughput, hasThroughput := d.GetOk("throughput"); hasThroughput {
		if throughput != 0 {
			table.CassandraTableCreateUpdateProperties.Options.Throughput = common.ConvertThroughputFromResourceDataLegacy(throughput)
		}
	}

	if _, hasAutoscaleSettings := d.GetOk("autoscale_settings"); hasAutoscaleSettings {
		table.CassandraTableCreateUpdateProperties.Options.AutoscaleSettings = common.ExpandCosmosDbAutoscaleSettingsLegacy(d)
	}

	future, err := client.CreateUpdateCassandraTable(ctx, resourceGroup, account, keyspace, name, table)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosDbCassandraTableRead(d, meta)
}

func resourceCosmosDbCassandraTableUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CassandraTableID(d.Id())
	if err != nil {
		return err
	}

	err = common.CheckForChangeFromAutoscaleAndManualThroughput(d)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	table := documentdb.CassandraTableCreateUpdateParameters{
		CassandraTableCreateUpdateProperties: &documentdb.CassandraTableCreateUpdateProperties{
			Resource: &documentdb.CassandraTableResource{
				ID: &id.TableName,
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	table.CassandraTableCreateUpdateProperties.Resource.Schema = expandTableSchema(d)

	if defaultTTL, hasTTL := d.GetOk("default_ttl"); hasTTL {
		table.CassandraTableCreateUpdateProperties.Resource.DefaultTTL = utils.Int32(int32(defaultTTL.(int)))
	}

	future, err := client.CreateUpdateCassandraTable(ctx, id.ResourceGroup, id.DatabaseAccountName, id.CassandraKeyspaceName, id.TableName, table)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", *id, err)
	}

	if common.HasThroughputChange(d) {
		throughputParameters := common.ExpandCosmosDBThroughputSettingsUpdateParametersLegacy(d)
		throughputFuture, err := client.UpdateCassandraTableThroughput(ctx, id.ResourceGroup, id.DatabaseAccountName, id.CassandraKeyspaceName, id.TableName, *throughputParameters)
		if err != nil {
			if response.WasNotFound(throughputFuture.Response()) {
				return fmt.Errorf("setting Throughput for %s: %+v - "+
					"If the collection has not been created with an initial throughput, you cannot configure it later", *id, err)
			}
		}

		if err = throughputFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("updating Throughput for %s: %+v", *id, err)
		}
	}

	return resourceCosmosDbCassandraTableRead(d, meta)
}

func resourceCosmosDbCassandraTableRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClient
	accountClient := meta.(*clients.Client).Cosmos.DatabaseClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	id, err := parse.CassandraTableID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetCassandraTable(ctx, id.ResourceGroup, id.DatabaseAccountName, id.CassandraKeyspaceName, id.TableName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	keyspaceId := parse.NewCassandraKeyspaceID(subscriptionId, id.ResourceGroup, id.DatabaseAccountName, id.CassandraKeyspaceName)

	d.Set("cassandra_keyspace_id", keyspaceId.ID())
	if props := resp.CassandraTableGetProperties; props != nil {
		if res := props.Resource; res != nil {
			d.Set("name", res.ID)

			if defaultTTL := res.DefaultTTL; defaultTTL != nil {
				d.Set("default_ttl", defaultTTL)
			}

			var analyticalTTL int
			if res.AnalyticalStorageTTL != nil {
				analyticalTTL = int(*res.AnalyticalStorageTTL)
			}
			d.Set("analytical_storage_ttl", analyticalTTL)

			if schema := res.Schema; schema != nil {
				d.Set("schema", flattenTableSchema(schema))
			}
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
		throughputResp, err := client.GetCassandraTableThroughput(ctx, id.ResourceGroup, id.DatabaseAccountName, id.CassandraKeyspaceName, id.TableName)
		if err != nil {
			if !utils.ResponseWasNotFound(throughputResp.Response) {
				return fmt.Errorf("retrieving Throughput for %s: %+v", *id, err)
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

func resourceCosmosDbCassandraTableDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CassandraTableID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.DeleteCassandraTable(ctx, id.ResourceGroup, id.DatabaseAccountName, id.CassandraKeyspaceName, id.TableName)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandTableSchema(d *pluginsdk.ResourceData) *documentdb.CassandraSchema {
	i := d.Get("schema").([]interface{})

	if len(i) == 0 || i[0] == nil {
		return nil
	}
	input := i[0].(map[string]interface{})

	cassandraSchema := documentdb.CassandraSchema{}

	if v, ok := input["column"].([]interface{}); ok {
		cassandraSchema.Columns = expandTableSchemaColumns(v)
	}

	if v, ok := input["partition_key"].([]interface{}); ok {
		cassandraSchema.PartitionKeys = expandTableSchemaPartitionKeys(v)
	}

	if v, ok := input["cluster_key"].([]interface{}); ok {
		cassandraSchema.ClusterKeys = expandTableSchemaClusterKeys(v)
	}

	return &cassandraSchema
}

func expandTableSchemaColumns(input []interface{}) *[]documentdb.Column {
	columns := make([]documentdb.Column, 0)
	for _, col := range input {
		data := col.(map[string]interface{})
		column := documentdb.Column{
			Name: utils.String(data["name"].(string)),
			Type: utils.String(data["type"].(string)),
		}
		columns = append(columns, column)
	}

	return &columns
}

func expandTableSchemaPartitionKeys(input []interface{}) *[]documentdb.CassandraPartitionKey {
	keys := make([]documentdb.CassandraPartitionKey, 0)
	for _, key := range input {
		data := key.(map[string]interface{})
		k := documentdb.CassandraPartitionKey{
			Name: utils.String(data["name"].(string)),
		}
		keys = append(keys, k)
	}

	return &keys
}

func expandTableSchemaClusterKeys(input []interface{}) *[]documentdb.ClusterKey {
	keys := make([]documentdb.ClusterKey, 0)
	for _, key := range input {
		data := key.(map[string]interface{})
		k := documentdb.ClusterKey{
			Name:    utils.String(data["name"].(string)),
			OrderBy: utils.String(data["order_by"].(string)),
		}
		keys = append(keys, k)
	}

	return &keys
}

func flattenTableSchema(input *documentdb.CassandraSchema) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})
	result["column"] = flattenTableSchemaColumns(input.Columns)
	result["partition_key"] = flattenTableSchemaPartitionKeys(input.PartitionKeys)
	result["cluster_key"] = flattenTableSchemaClusterKeys(input.ClusterKeys)

	results = append(results, result)
	return results
}

func flattenTableSchemaColumns(input *[]documentdb.Column) []interface{} {
	if input == nil {
		return nil
	}

	columns := make([]interface{}, 0)

	for _, v := range *input {
		name := ""
		if v.Name != nil {
			name = *v.Name
		}
		typeStr := ""
		if v.Type != nil {
			typeStr = *v.Type
		}
		columns = append(columns, map[string]interface{}{
			"name": name,
			"type": typeStr,
		})
	}

	return columns
}

func flattenTableSchemaPartitionKeys(input *[]documentdb.CassandraPartitionKey) []interface{} {
	if input == nil {
		return nil
	}

	keys := make([]interface{}, 0)

	for _, v := range *input {
		name := ""
		if v.Name != nil {
			name = *v.Name
		}
		keys = append(keys, map[string]interface{}{
			"name": name,
		})
	}

	return keys
}

func flattenTableSchemaClusterKeys(input *[]documentdb.ClusterKey) []interface{} {
	if input == nil {
		return nil
	}

	keys := make([]interface{}, 0)

	for _, v := range *input {
		name := ""
		if v.Name != nil {
			name = *v.Name
		}
		orderBy := ""
		if v.OrderBy != nil {
			orderBy = *v.OrderBy
		}
		keys = append(keys, map[string]interface{}{
			"name":     name,
			"order_by": orderBy,
		})
	}

	return keys
}
