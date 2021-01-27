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

	existing, err := client.GetCassandraTable(ctx, resourceGroup, account, keyspace, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of creating Cosmos Cassandra Table %q (Account: %q, Keyspace: %q): %+v", name, account, keyspace, err)
		}
	} else {
		if existing.ID == nil && *existing.ID == "" {
			return fmt.Errorf("Error generating import ID for Cosmos Cassandra Table %q (Account: %q, Keyspace: %q)", name, account, keyspace)
		}

		return tf.ImportAsExistsError("azurerm_cosmosdb_cassandra_table", *existing.ID)
	}

	table := documentdb.CassandraTableCreateUpdateParameters{
		CassandraTableCreateUpdateProperties: &documentdb.CassandraTableCreateUpdateProperties{
			Resource: &documentdb.CassandraTableResource{
				ID: &name,
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

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
		return fmt.Errorf("Error issuing create/update request for Cosmos Cassandra Table %q (Account: %q, Keyspace: %q): %+v", name, account, keyspace, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Cosmos Cassandra Table %q (Account: %q, Keyspace: %q): %+v", name, account, keyspace, err)
	}

	resp, err := client.GetCassandraTable(ctx, resourceGroup, account, keyspace, name)
	if err != nil {
		return fmt.Errorf("Error making get request for Cosmos Cassandra Table %q (Account: %q, Keyspace: %q): %+v", name, account, keyspace, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Error getting ID from Cosmos Cassandra Table %q (Account: %q, Keyspace: %q)", name, account, keyspace)
	}

	d.SetId(*resp.ID)

	return resourceCosmosDbCassandraTableRead(d, meta)
}

func resourceCosmosDbCassandraTableUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CassandraTableID(d.Id())
	if err != nil {
		return err
	}

	err = common.CheckForChangeFromAutoscaleAndManualThroughput(d)
	if err != nil {
		return fmt.Errorf("Error updating Cosmos Cassandra Table %q (Account: %q, Keyspace: %q) - %+v", id.TableName, id.DatabaseAccountName, id.CassandraKeyspaceName, err)
	}

	table := documentdb.CassandraTableCreateUpdateParameters{
		CassandraTableCreateUpdateProperties: &documentdb.CassandraTableCreateUpdateProperties{
			Resource: &documentdb.CassandraTableResource{
				ID: &id.TableName,
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	if defaultTTL, hasTTL := d.GetOk("default_ttl"); hasTTL {
		table.CassandraTableCreateUpdateProperties.Resource.DefaultTTL = utils.Int32(int32(defaultTTL.(int)))
	}

	future, err := client.CreateUpdateCassandraTable(ctx, id.ResourceGroup, id.DatabaseAccountName, id.DatabaseAccountName, id.TableName, table)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Cosmos Cassandra Table %q (Account: %q, Keyspace: %q): %+v", id.ResourceGroup, id.DatabaseAccountName, id.CassandraKeyspaceName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Cosmos Cassandra Table %q (Account: %q, Keyspace: %q): %+v", id.ResourceGroup, id.DatabaseAccountName, id.CassandraKeyspaceName, err)
	}

	if common.HasThroughputChange(d) {
		throughputParameters := common.ExpandCosmosDBThroughputSettingsUpdateParameters(d)
		throughputFuture, err := client.UpdateCassandraTableThroughput(ctx, id.ResourceGroup, id.DatabaseAccountName, id.CassandraKeyspaceName, id.TableName, *throughputParameters)
		if err != nil {
			if response.WasNotFound(throughputFuture.Response()) {
				return fmt.Errorf("Error setting Throughput for Cosmos Cassandra Table %q (Account: %q, Keyspace: %q): %+v - "+
					"If the collection has not been created with an initial throughput, you cannot configure it later.", id.TableName, id.DatabaseAccountName, id.CassandraKeyspaceName, err)
			}
		}

		if err = throughputFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting on ThroughputUpdate future for Cosmos Cassandra Table %q (Account: %q, Keyspace: %q): %+v", id.TableName, id.DatabaseAccountName, id.CassandraKeyspaceName, err)
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
			log.Printf("[INFO] Error reading Cosmos Cassandra Table %q (Account: %q, Keyspace: %q) - removing from state", id.TableName, id.DatabaseAccountName, id.CassandraKeyspaceName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Cosmos Cassandra Table %q (Account: %q, Keyspace: %q): %+v", id.TableName, id.DatabaseAccountName, id.CassandraKeyspaceName, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_name", id.DatabaseAccountName)
	if props := resp.CassandraTableGetProperties; props != nil {
		if res := props.Resource; res != nil {
			d.Set("name", res.ID)

			if defaultTTL := res.DefaultTTL; defaultTTL != nil {
				d.Set("default_ttl", defaultTTL)
			}
		}
	}

	throughputResp, err := client.GetCassandraTableThroughput(ctx, id.ResourceGroup, id.DatabaseAccountName, id.CassandraKeyspaceName, id.TableName)
	if err != nil {
		if !utils.ResponseWasNotFound(throughputResp.Response) {
			return fmt.Errorf("Error reading Throughput on Cosmos Cassandra Table %q (Account: %q, Keyspace: %q): %+v", id.TableName, id.DatabaseAccountName, id.CassandraKeyspaceName, err)
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
			return fmt.Errorf("Error deleting Cosmos Cassandra Table %q (Account: %q, Keyspace: %q): %+v", id.TableName, id.DatabaseAccountName, id.CassandraKeyspaceName, err)
		}
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting on delete future for Cosmos Cassandra Table %q (Account: %q, Keyspace: %q): %+v", id.TableName, id.DatabaseAccountName, id.CassandraKeyspaceName, err)
	}

	return nil
}
