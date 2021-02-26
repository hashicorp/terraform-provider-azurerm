package cosmos

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/cosmos-db/mgmt/2020-04-01-preview/documentdb"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/migration"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceCosmosDbSQLDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmosDbSQLDatabaseCreate,
		Read:   resourceCosmosDbSQLDatabaseRead,
		Update: resourceCosmosDbSQLDatabaseUpdate,
		Delete: resourceCosmosDbSQLDatabaseDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    migration.ResourceSqlDatabaseUpgradeV0Schema().CoreConfigSchema().ImpliedType(),
				Upgrade: migration.ResourceSqlDatabaseStateUpgradeV0ToV1,
				Version: 0,
			},
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

func resourceCosmosDbSQLDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	account := d.Get("account_name").(string)

	existing, err := client.GetSQLDatabase(ctx, resourceGroup, account, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of creating Cosmos SQL Database %q (Account: %q): %+v", name, account, err)
		}
	} else {
		if existing.ID == nil && *existing.ID == "" {
			return fmt.Errorf("Error generating import ID for Cosmos SQL Database %q (Account: %q)", name, account)
		}

		return tf.ImportAsExistsError("azurerm_cosmosdb_sql_database", *existing.ID)
	}

	db := documentdb.SQLDatabaseCreateUpdateParameters{
		SQLDatabaseCreateUpdateProperties: &documentdb.SQLDatabaseCreateUpdateProperties{
			Resource: &documentdb.SQLDatabaseResource{
				ID: &name,
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	if throughput, hasThroughput := d.GetOk("throughput"); hasThroughput {
		if throughput != 0 {
			db.SQLDatabaseCreateUpdateProperties.Options.Throughput = common.ConvertThroughputFromResourceData(throughput)
		}
	}

	if _, hasAutoscaleSettings := d.GetOk("autoscale_settings"); hasAutoscaleSettings {
		db.SQLDatabaseCreateUpdateProperties.Options.AutoscaleSettings = common.ExpandCosmosDbAutoscaleSettings(d)
	}

	future, err := client.CreateUpdateSQLDatabase(ctx, resourceGroup, account, name, db)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Cosmos SQL Database %q (Account: %q): %+v", name, account, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Cosmos SQL Database %q (Account: %q): %+v", name, account, err)
	}

	resp, err := client.GetSQLDatabase(ctx, resourceGroup, account, name)
	if err != nil {
		return fmt.Errorf("Error making get request for Cosmos SQL Database %q (Account: %q): %+v", name, account, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Error getting ID from Cosmos SQL Database %q (Account: %q)", name, account)
	}

	d.SetId(*resp.ID)

	return resourceCosmosDbSQLDatabaseRead(d, meta)
}

func resourceCosmosDbSQLDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlDatabaseID(d.Id())
	if err != nil {
		return err
	}

	err = common.CheckForChangeFromAutoscaleAndManualThroughput(d)
	if err != nil {
		return fmt.Errorf("Error updating Cosmos SQL Database %q (Account: %q) - %+v", id.Name, id.DatabaseAccountName, err)
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
		return fmt.Errorf("Error issuing create/update request for Cosmos SQL Database %q (Account: %q): %+v", id.Name, id.DatabaseAccountName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Cosmos SQL Database %q (Account: %q): %+v", id.Name, id.DatabaseAccountName, err)
	}

	if common.HasThroughputChange(d) {
		throughputParameters := common.ExpandCosmosDBThroughputSettingsUpdateParameters(d)
		throughputFuture, err := client.UpdateSQLDatabaseThroughput(ctx, id.ResourceGroup, id.DatabaseAccountName, id.Name, *throughputParameters)
		if err != nil {
			if response.WasNotFound(throughputFuture.Response()) {
				return fmt.Errorf("Error setting Throughput for Cosmos SQL Database %q (Account: %q) %+v - "+
					"If the collection has not been created with an initial throughput, you cannot configure it later.", id.Name, id.DatabaseAccountName, err)
			}
		}

		if err = throughputFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting on ThroughputUpdate future for Cosmos SQL Database %q (Account: %q): %+v", id.Name, id.DatabaseAccountName, err)
		}
	}

	return resourceCosmosDbSQLDatabaseRead(d, meta)
}

func resourceCosmosDbSQLDatabaseRead(d *schema.ResourceData, meta interface{}) error {
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

		return fmt.Errorf("Error reading Cosmos SQL Database %q (Account: %q): %+v", id.Name, id.DatabaseAccountName, err)
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

	if props := accResp.DatabaseAccountGetProperties; props != nil && props.Capabilities != nil {
		serverless := false
		for _, v := range *props.Capabilities {
			if *v.Name == "EnableServerless" {
				serverless = true
			}
		}

		if !serverless {
			throughputResp, err := client.GetSQLDatabaseThroughput(ctx, id.ResourceGroup, id.DatabaseAccountName, id.Name)
			if err != nil {
				if !utils.ResponseWasNotFound(throughputResp.Response) {
					return fmt.Errorf("Error reading Throughput on Cosmos SQL Database %q (Account: %q) ID: %v", id.Name, id.DatabaseAccountName, err)
				} else {
					d.Set("throughput", nil)
					d.Set("autoscale_settings", nil)
				}
			} else {
				common.SetResourceDataThroughputFromResponse(throughputResp, d)
			}
		}
	}

	return nil
}

func resourceCosmosDbSQLDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
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
			return fmt.Errorf("Error deleting Cosmos SQL Database %q (Account: %q): %+v", id.Name, id.DatabaseAccountName, err)
		}
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting on delete future for Cosmos SQL Database %q (Account: %q): %+v", id.Name, id.DatabaseAccountName, err)
	}

	return nil
}
