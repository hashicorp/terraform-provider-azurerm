// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"fmt"
	"math"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbCassandraTable() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCosmosDbCassandraTableCreate,
		Read:   resourceCosmosDbCassandraTableRead,
		Update: resourceCosmosDbCassandraTableUpdate,
		Delete: resourceCosmosDbCassandraTableDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := cosmosdb.ParseCassandraKeyspaceTableID(id)
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
				ValidateFunc: cosmosdb.ValidateCassandraKeyspaceID,
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
					validation.IntBetween(-1, math.MaxInt32),
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
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	keyspaceId, err := cosmosdb.ParseCassandraKeyspaceID(d.Get("cassandra_keyspace_id").(string))
	if err != nil {
		return err
	}

	id := cosmosdb.NewCassandraKeyspaceTableID(meta.(*clients.Client).Account.SubscriptionId, keyspaceId.ResourceGroupName, keyspaceId.DatabaseAccountName, keyspaceId.CassandraKeyspaceName, d.Get("name").(string))

	existing, err := client.CassandraResourcesGetCassandraTable(ctx, id)
	if !response.WasNotFound(existing.HttpResponse) {
		if err != nil {
			return fmt.Errorf("checking for presence of existing %+v: %+v", id, err)
		}
		return tf.ImportAsExistsError("azurerm_cosmosdb_cassandra_table", id.ID())
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

	if err := client.CassandraResourcesCreateUpdateCassandraTableThenPoll(ctx, id, table); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosDbCassandraTableRead(d, meta)
}

func resourceCosmosDbCassandraTableUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceCosmosDbCassandraTableRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

	if !isServerlessCapacityMode(accResp.Model) {
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

func resourceCosmosDbCassandraTableDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseCassandraKeyspaceTableID(d.Id())
	if err != nil {
		return err
	}

	if err := client.CassandraResourcesDeleteCassandraTableThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandTableSchema(d *pluginsdk.ResourceData) *cosmosdb.CassandraSchema {
	i := d.Get("schema").([]interface{})

	if len(i) == 0 || i[0] == nil {
		return nil
	}
	input := i[0].(map[string]interface{})

	cassandraSchema := cosmosdb.CassandraSchema{}

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

func expandTableSchemaColumns(input []interface{}) *[]cosmosdb.Column {
	columns := make([]cosmosdb.Column, 0)
	for _, col := range input {
		data := col.(map[string]interface{})
		column := cosmosdb.Column{
			Name: pointer.To(data["name"].(string)),
			Type: pointer.To(data["type"].(string)),
		}
		columns = append(columns, column)
	}

	return &columns
}

func expandTableSchemaPartitionKeys(input []interface{}) *[]cosmosdb.CassandraPartitionKey {
	keys := make([]cosmosdb.CassandraPartitionKey, 0)
	for _, key := range input {
		data := key.(map[string]interface{})
		k := cosmosdb.CassandraPartitionKey{
			Name: pointer.To(data["name"].(string)),
		}
		keys = append(keys, k)
	}

	return &keys
}

func expandTableSchemaClusterKeys(input []interface{}) *[]cosmosdb.ClusterKey {
	keys := make([]cosmosdb.ClusterKey, 0)
	for _, key := range input {
		data := key.(map[string]interface{})
		k := cosmosdb.ClusterKey{
			Name:    pointer.To(data["name"].(string)),
			OrderBy: pointer.To(data["order_by"].(string)),
		}
		keys = append(keys, k)
	}

	return &keys
}

func flattenTableSchema(input *cosmosdb.CassandraSchema) []interface{} {
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

func flattenTableSchemaColumns(input *[]cosmosdb.Column) []interface{} {
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

func flattenTableSchemaPartitionKeys(input *[]cosmosdb.CassandraPartitionKey) []interface{} {
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

func flattenTableSchemaClusterKeys(input *[]cosmosdb.ClusterKey) []interface{} {
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
