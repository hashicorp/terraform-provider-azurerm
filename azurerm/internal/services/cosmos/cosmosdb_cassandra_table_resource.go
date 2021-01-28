package cosmos

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/Azure/azure-sdk-for-go/services/preview/cosmos-db/mgmt/2020-04-01-preview/documentdb"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceCosmosDbCassandraTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmosDbCassandraTableCreate,
		Read:   resourceCosmosDbCassandraTableRead,
		Update: resourceCosmosDbCassandraTableUpdate,
		Delete: resourceCosmosDbCassandraTableDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosEntityName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosAccountName,
			},

			"default_ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(-1),
			},

			"keyspace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosEntityName,
			},

			"schema": common.CassandraTableSchemaPropertySchema(),

			"throughput": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.CosmosThroughput,
			},

			"autoscale_settings": common.DatabaseAutoscaleSettingsSchema(),
		},
	}
}

func resourceCosmosDbCassandraTableCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	account := d.Get("account_name").(string)
	keyspace := d.Get("keyspace_name").(string)

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := parse.NewCassandraTableID(subscriptionId, resourceGroup, account, keyspace, name)
	existing, err := client.GetCassandraTable(ctx, resourceGroup, account, keyspace, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %+v: %+v", id, err)
		}
	} else {
		if existing.ID == nil && *existing.ID == "" {
			return fmt.Errorf("Error generating import ID for Cosmos Cassandra Table %q (Account: %q, Keyspace: %q)", name, account, keyspace)
		}

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

	if throughput, hasThroughput := d.GetOk("throughput"); hasThroughput {
		if throughput != 0 {
			table.CassandraTableCreateUpdateProperties.Options.Throughput = common.ConvertThroughputFromResourceData(throughput)
		}
	}

	if _, hasAutoscaleSettings := d.GetOk("autoscale_settings"); hasAutoscaleSettings {
		table.CassandraTableCreateUpdateProperties.Options.AutoscaleSettings = common.ExpandCosmosDbAutoscaleSettings(d)
	}

	future, err := client.CreateUpdateCassandraTable(ctx, resourceGroup, account, keyspace, name, table)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	resp, err := client.GetCassandraTable(ctx, resourceGroup, account, keyspace, name)
	if err != nil {
		return fmt.Errorf("Error making get request for Cosmos Cassandra Table %q (Account: %q, Keyspace: %q): %+v", name, account, keyspace, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Error getting ID from Cosmos Cassandra Table %q (Account: %q, Keyspace: %q)", name, account, keyspace)
	}

	d.SetId(id.ID())

	return resourceCosmosDbCassandraTableRead(d, meta)
}

func resourceCosmosDbCassandraTableUpdate(d *schema.ResourceData, meta interface{}) error {
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

	future, err := client.CreateUpdateCassandraTable(ctx, id.ResourceGroup, id.DatabaseAccountName, id.DatabaseAccountName, id.TableName, table)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Cosmos Cassandra Table %q (Account: %q, Keyspace: %q): %+v", id.ResourceGroup, id.DatabaseAccountName, id.CassandraKeyspaceName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for update of %s: %+v", *id, err)
	}

	if common.HasThroughputChange(d) {
		throughputParameters := common.ExpandCosmosDBThroughputSettingsUpdateParameters(d)
		throughputFuture, err := client.UpdateCassandraTableThroughput(ctx, id.ResourceGroup, id.DatabaseAccountName, id.CassandraKeyspaceName, id.TableName, *throughputParameters)
		if err != nil {
			if response.WasNotFound(throughputFuture.Response()) {
				return fmt.Errorf("setting Throughput for %s: %+v - "+
					"If the collection has not been created with an initial throughput, you cannot configure it later.", *id, err)
			}
		}

		if err = throughputFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("updating Throughput for %s: %+v", *id, err)
		}
	}

	return resourceCosmosDbCassandraTableRead(d, meta)
}

func resourceCosmosDbCassandraTableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
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

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_name", id.DatabaseAccountName)
	d.Set("keyspace_name", id.CassandraKeyspaceName)
	if props := resp.CassandraTableGetProperties; props != nil {
		if res := props.Resource; res != nil {
			d.Set("name", res.ID)

			if defaultTTL := res.DefaultTTL; defaultTTL != nil {
				d.Set("default_ttl", defaultTTL)
			}

			if schema := res.Schema; schema != nil {
				d.Set("schema", flattenTableSchema(schema))
			}
		}
	}

	throughputResp, err := client.GetCassandraTableThroughput(ctx, id.ResourceGroup, id.DatabaseAccountName, id.CassandraKeyspaceName, id.TableName)
	if err != nil {
		if !utils.ResponseWasNotFound(throughputResp.Response) {
			return fmt.Errorf("retrieving Throughput for %s: %+v", *id, err)
		} else {
			d.Set("throughput", nil)
			d.Set("autoscale_settings", nil)
		}
	} else {
		common.SetResourceDataThroughputFromResponse(throughputResp, d)
	}

	return nil
}

func resourceCosmosDbCassandraTableDelete(d *schema.ResourceData, meta interface{}) error {
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

func expandTableSchema(d *schema.ResourceData) *documentdb.CassandraSchema {
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

func flattenTableSchema(schema *documentdb.CassandraSchema) []interface{} {
	results := make([]interface{}, 0)
	if schema == nil {
		return results
	}

	result := make(map[string]interface{})
	result["column"] = flattenTableSchemaColumns(schema.Columns)
	result["partition_key"] = flattenTableSchemaPartitionKeys(schema.PartitionKeys)
	result["cluster_key"] = flattenTableSchemaClusterKeys(schema.ClusterKeys)

	results = append(results, result)
	return results
}

func flattenTableSchemaColumns(input *[]documentdb.Column) []interface{} {
	if input == nil {
		return nil
	}

	columns := make([]interface{}, 0)

	for _, v := range *input {
		block := make(map[string]interface{})
		block["name"] = v.Name
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
		columns = append(columns, block)
	}

	return columns
}

func flattenTableSchemaPartitionKeys(input *[]documentdb.CassandraPartitionKey) []interface{} {
	if input == nil {
		return nil
	}

	keys := make([]interface{}, 0)

	for _, v := range *input {
		block := make(map[string]interface{})
		block["name"] = v.Name
		keys = append(keys, block)
	}

	return keys
}

func flattenTableSchemaClusterKeys(input *[]documentdb.ClusterKey) []interface{} {
	if input == nil {
		return nil
	}

	keys := make([]interface{}, 0)

	for _, v := range *input {
		block := make(map[string]interface{})
		block["name"] = v.Name
		block["order_by"] = v.OrderBy
		keys = append(keys, block)
	}

	return keys
}
