package cosmos

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmCosmosGremlinDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCosmosGremlinDatabaseCreate,
		Update: resourceArmCosmosGremlinDatabaseUpdate,
		Read:   resourceArmCosmosGremlinDatabaseRead,
		Delete: resourceArmCosmosGremlinDatabaseDelete,

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

			"throughput": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.CosmosThroughput,
			},
		},
	}
}

func resourceArmCosmosGremlinDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.DatabaseClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	account := d.Get("account_name").(string)

	if features.ShouldResourcesBeImported() {
		existing, err := client.GetGremlinDatabase(ctx, resourceGroup, account, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of creating Gremlin Database %s (Account %s): %+v", name, account, err)
			}
		} else {
			id, err := azure.CosmosGetIDFromResponse(existing.Response)
			if err != nil {
				return fmt.Errorf("Error generating import ID for Cosmos Gremlin Database '%s' (Account %s)", name, account)
			}

			return tf.ImportAsExistsError("azurerm_cosmosdb_gremlin_database", id)
		}
	}

	db := documentdb.GremlinDatabaseCreateUpdateParameters{
		GremlinDatabaseCreateUpdateProperties: &documentdb.GremlinDatabaseCreateUpdateProperties{
			Resource: &documentdb.GremlinDatabaseResource{
				ID: &name,
			},
			Options: map[string]*string{},
		},
	}

	if throughput, hasThroughput := d.GetOk("throughput"); hasThroughput {
		db.GremlinDatabaseCreateUpdateProperties.Options = map[string]*string{
			"throughput": utils.String(strconv.Itoa(throughput.(int))),
		}
	}

	future, err := client.CreateUpdateGremlinDatabase(ctx, resourceGroup, account, name, db)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Gremlin Database %s (Account %s): %+v", name, account, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Cosmos Gremlin Database %s (Account %s): %+v", name, account, err)
	}

	resp, err := client.GetGremlinDatabase(ctx, resourceGroup, account, name)
	if err != nil {
		return fmt.Errorf("Error making get request for Cosmos Gremlin Database %s (Account %s) ID: %v", name, account, err)
	}

	id, err := azure.CosmosGetIDFromResponse(resp.Response)
	if err != nil {
		return fmt.Errorf("Error retrieving the ID for Cosmos Gremlin Database '%s' (Account %s) ID: %v", name, account, err)
	}
	d.SetId(id)

	return resourceArmCosmosGremlinDatabaseRead(d, meta)
}

func resourceArmCosmosGremlinDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.DatabaseClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseCosmosDatabaseID(d.Id())
	if err != nil {
		return err
	}

	db := documentdb.GremlinDatabaseCreateUpdateParameters{
		GremlinDatabaseCreateUpdateProperties: &documentdb.GremlinDatabaseCreateUpdateProperties{
			Resource: &documentdb.GremlinDatabaseResource{
				ID: &id.Database,
			},
			Options: map[string]*string{},
		},
	}

	future, err := client.CreateUpdateGremlinDatabase(ctx, id.ResourceGroup, id.Account, id.Database, db)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Cosmos Gremlin Database %s (Account %s): %+v", id.Database, id.Account, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Cosmos Gremlin Database %s (Account %s): %+v", id.Database, id.Account, err)
	}

	if d.HasChange("throughput") {
		throughputParameters := documentdb.ThroughputUpdateParameters{
			ThroughputUpdateProperties: &documentdb.ThroughputUpdateProperties{
				Resource: &documentdb.ThroughputResource{
					Throughput: utils.Int32(int32(d.Get("throughput").(int))),
				},
			},
		}

		throughputFuture, err := client.UpdateGremlinDatabaseThroughput(ctx, id.ResourceGroup, id.Account, id.Database, throughputParameters)
		if err != nil {
			if response.WasNotFound(throughputFuture.Response()) {
				return fmt.Errorf("Error setting Throughput for Cosmos Gremlin Database %s (Account %s): %+v - "+
					"If the collection has not been created with and initial throughput, you cannot configure it later.", id.Database, id.Account, err)
			}
		}

		if err = throughputFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting on ThroughputUpdate future for Cosmos Gremlin Database %s (Account %s, Database %s): %+v", id.Database, id.Account, id.Database, err)
		}
	}

	if _, err = client.GetGremlinDatabase(ctx, id.ResourceGroup, id.Account, id.Database); err != nil {
		return fmt.Errorf("Error making get request for Cosmos Gremlin Database %s (Account %s): %+v", id.Database, id.Account, err)
	}

	return resourceArmCosmosGremlinDatabaseRead(d, meta)
}

func resourceArmCosmosGremlinDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.DatabaseClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseCosmosDatabaseID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetGremlinDatabase(ctx, id.ResourceGroup, id.Account, id.Database)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Cosmos Gremlin Database %s (Account %s) - removing from state", id.Database, id.Account)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Cosmos Gremlin Database %s (Account %s): %+v", id.Database, id.Account, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_name", id.Account)
	if props := resp.GremlinDatabaseProperties; props != nil {
		d.Set("name", props.ID)
	}

	throughputResp, err := client.GetGremlinDatabaseThroughput(ctx, id.ResourceGroup, id.Account, id.Database)
	if err != nil {
		if !utils.ResponseWasNotFound(throughputResp.Response) {
			return fmt.Errorf("Error reading Throughput on Cosmos Gremlin Database %s (Account %s): %+v", id.Database, id.Account, err)
		} else {
			d.Set("throughput", nil)
		}
	} else {
		d.Set("throughput", throughputResp.Throughput)
	}

	return nil
}

func resourceArmCosmosGremlinDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.DatabaseClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseCosmosDatabaseID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.DeleteGremlinDatabase(ctx, id.ResourceGroup, id.Account, id.Database)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting Cosmos Gremlin Database %s (Account %s): %+v", id.Database, id.Account, err)
		}
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on delete future for Cosmos Gremlin Database %s (Account %s): %+v", id.Database, id.Account, err)
	}

	return nil
}
