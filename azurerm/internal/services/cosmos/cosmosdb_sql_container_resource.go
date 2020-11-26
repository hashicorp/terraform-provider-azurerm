package cosmos

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/cosmos-db/mgmt/2020-04-01-preview/documentdb"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/migration"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmCosmosDbSQLContainer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCosmosDbSQLContainerCreate,
		Read:   resourceArmCosmosDbSQLContainerRead,
		Update: resourceArmCosmosDbSQLContainerUpdate,
		Delete: resourceArmCosmosDbSQLContainerDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    migration.ResourceSqlContainerUpgradeV0Schema().CoreConfigSchema().ImpliedType(),
				Upgrade: migration.ResourceSqlContainerStateUpgradeV0ToV1,
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

			"database_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosEntityName,
			},

			"partition_key_path": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"partition_key_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 2),
			},

			"throughput": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.CosmosThroughput,
			},

			"autoscale_settings": common.ContainerAutoscaleSettingsSchema(),

			"default_ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(-1),
			},

			"unique_key": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"paths": {
							Type:     schema.TypeSet,
							Required: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},
			"indexing_policy": common.CosmosDbIndexingPolicySchema(),
		},
	}
}

func resourceArmCosmosDbSQLContainerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	database := d.Get("database_name").(string)
	account := d.Get("account_name").(string)
	partitionkeypaths := d.Get("partition_key_path").(string)
	partitionKeyVersion, hasPartitionKeyVersion := d.GetOk("partition_key_version")

	existing, err := client.GetSQLContainer(ctx, resourceGroup, account, database, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of creating Cosmos SQL Container %q (Account: %q, Database: %q): %+v", name, account, database, err)
		}
	} else {
		if existing.ID == nil && *existing.ID == "" {
			return fmt.Errorf("Error generating import ID for Cosmos SQL Container %q (Account: %q, Database: %q)", name, account, database)
		}

		return tf.ImportAsExistsError("azurerm_cosmosdb_sql_container", *existing.ID)
	}

	indexingPolicy := common.ExpandAzureRmCosmosDbIndexingPolicy(d)
	err = common.ValidateAzureRmCosmosDbIndexingPolicy(indexingPolicy)
	if err != nil {
		return fmt.Errorf("Error generating indexing policy for Cosmos SQL Container %q (Account: %q, Database: %q)", name, account, database)
	}

	db := documentdb.SQLContainerCreateUpdateParameters{
		SQLContainerCreateUpdateProperties: &documentdb.SQLContainerCreateUpdateProperties{
			Resource: &documentdb.SQLContainerResource{
				ID:             &name,
				IndexingPolicy: indexingPolicy,
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	if partitionkeypaths != "" {
		if hasPartitionKeyVersion {
			db.SQLContainerCreateUpdateProperties.Resource.PartitionKey = &documentdb.ContainerPartitionKey{
				Paths:   &[]string{partitionkeypaths},
				Kind:    documentdb.PartitionKindHash,
				Version: utils.Int32(int32(partitionKeyVersion.(int))),
			}
		} else {
			db.SQLContainerCreateUpdateProperties.Resource.PartitionKey = &documentdb.ContainerPartitionKey{
				Paths: &[]string{partitionkeypaths},
				Kind:  documentdb.PartitionKindHash,
			}
		}
	}

	if keys := expandCosmosSQLContainerUniqueKeys(d.Get("unique_key").(*schema.Set)); keys != nil {
		db.SQLContainerCreateUpdateProperties.Resource.UniqueKeyPolicy = &documentdb.UniqueKeyPolicy{
			UniqueKeys: keys,
		}
	}

	if defaultTTL, hasTTL := d.GetOk("default_ttl"); hasTTL {
		db.SQLContainerCreateUpdateProperties.Resource.DefaultTTL = utils.Int32(int32(defaultTTL.(int)))
	}

	if throughput, hasThroughput := d.GetOk("throughput"); hasThroughput {
		if throughput != 0 {
			db.SQLContainerCreateUpdateProperties.Options.Throughput = common.ConvertThroughputFromResourceData(throughput)
		}
	}

	if _, hasAutoscaleSettings := d.GetOk("autoscale_settings"); hasAutoscaleSettings {
		db.SQLContainerCreateUpdateProperties.Options.AutoscaleSettings = common.ExpandCosmosDbAutoscaleSettings(d)
	}

	future, err := client.CreateUpdateSQLContainer(ctx, resourceGroup, account, database, name, db)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Cosmos SQL Container %q (Account: %q, Database: %q): %+v", name, account, database, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Cosmos SQL Container %q (Account: %q, Database: %q): %+v", name, account, database, err)
	}

	resp, err := client.GetSQLContainer(ctx, resourceGroup, account, database, name)
	if err != nil {
		return fmt.Errorf("Error making get request for Cosmos SQL Container %q (Account: %q, Database: %q): %+v", name, account, database, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Error getting ID from Cosmos SQL Container %q (Account: %q, Database: %q)", name, account, database)
	}

	d.SetId(*resp.ID)

	return resourceArmCosmosDbSQLContainerRead(d, meta)
}

func resourceArmCosmosDbSQLContainerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlContainerID(d.Id())
	if err != nil {
		return err
	}

	err = common.CheckForChangeFromAutoscaleAndManualThroughput(d)
	if err != nil {
		return fmt.Errorf("Error updating Cosmos SQL Container %q (Account: %q, Database: %q): %+v", id.ContainerName, id.DatabaseAccountName, id.SqlDatabaseName, err)
	}

	partitionkeypaths := d.Get("partition_key_path").(string)
	partitionKeyVersion, hasPartitionKeyVersion := d.GetOk("partition_key_version")

	indexingPolicy := common.ExpandAzureRmCosmosDbIndexingPolicy(d)
	err = common.ValidateAzureRmCosmosDbIndexingPolicy(indexingPolicy)
	if err != nil {
		return fmt.Errorf("Error updating Cosmos SQL Container %q (Account: %q, Database: %q): %+v", id.ContainerName, id.DatabaseAccountName, id.SqlDatabaseName, err)
	}

	db := documentdb.SQLContainerCreateUpdateParameters{
		SQLContainerCreateUpdateProperties: &documentdb.SQLContainerCreateUpdateProperties{
			Resource: &documentdb.SQLContainerResource{
				ID:             &id.ContainerName,
				IndexingPolicy: indexingPolicy,
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	if partitionkeypaths != "" {
		if hasPartitionKeyVersion {
			db.SQLContainerCreateUpdateProperties.Resource.PartitionKey = &documentdb.ContainerPartitionKey{
				Paths:   &[]string{partitionkeypaths},
				Kind:    documentdb.PartitionKindHash,
				Version: utils.Int32(int32(partitionKeyVersion.(int))),
			}
		} else {
			db.SQLContainerCreateUpdateProperties.Resource.PartitionKey = &documentdb.ContainerPartitionKey{
				Paths: &[]string{partitionkeypaths},
				Kind:  documentdb.PartitionKindHash,
			}
		}
	}

	if keys := expandCosmosSQLContainerUniqueKeys(d.Get("unique_key").(*schema.Set)); keys != nil {
		db.SQLContainerCreateUpdateProperties.Resource.UniqueKeyPolicy = &documentdb.UniqueKeyPolicy{
			UniqueKeys: keys,
		}
	}

	if defaultTTL, hasTTL := d.GetOk("default_ttl"); hasTTL {
		db.SQLContainerCreateUpdateProperties.Resource.DefaultTTL = utils.Int32(int32(defaultTTL.(int)))
	}

	future, err := client.CreateUpdateSQLContainer(ctx, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName, db)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Cosmos SQL Container %q (Account: %q, Database: %q): %+v", id.ContainerName, id.DatabaseAccountName, id.SqlDatabaseName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Cosmos SQL Container %q (Account: %q, Database: %q): %+v", id.ContainerName, id.DatabaseAccountName, id.SqlDatabaseName, err)
	}

	if common.HasThroughputChange(d) {
		throughputParameters := common.ExpandCosmosDBThroughputSettingsUpdateParameters(d)
		throughputFuture, err := client.UpdateSQLContainerThroughput(ctx, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName, *throughputParameters)
		if err != nil {
			if response.WasNotFound(throughputFuture.Response()) {
				return fmt.Errorf("Error setting Throughput for Cosmos SQL Container %q (Account: %q, Database: %q): %+v - "+
					"If the collection has not been created with an initial throughput, you cannot configure it later.", id.ContainerName, id.DatabaseAccountName, id.SqlDatabaseName, err)
			}
		}

		if err = throughputFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting on ThroughputUpdate future for Cosmos Container %q (Account: %q, Database: %q): %+v", id.ContainerName, id.DatabaseAccountName, id.SqlDatabaseName, err)
		}
	}

	return resourceArmCosmosDbSQLContainerRead(d, meta)
}

func resourceArmCosmosDbSQLContainerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlClient
	accountClient := meta.(*clients.Client).Cosmos.DatabaseClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlContainerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetSQLContainer(ctx, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Cosmos SQL Container %q (Account: %q) - removing from state", id.SqlDatabaseName, id.ContainerName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Cosmos SQL Container %q (Account: %q): %+v", id.SqlDatabaseName, id.ContainerName, err)
	}

	d.Set("name", id.ContainerName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_name", id.DatabaseAccountName)
	d.Set("database_name", id.SqlDatabaseName)

	if props := resp.SQLContainerGetProperties; props != nil {
		if res := props.Resource; res != nil {
			if pk := res.PartitionKey; pk != nil {
				if paths := pk.Paths; paths != nil {
					if len(*paths) > 1 {
						return fmt.Errorf("Error reading PartitionKey Paths, more then 1 returned")
					} else if len(*paths) == 1 {
						d.Set("partition_key_path", (*paths)[0])
					}
				}
				if version := pk.Version; version != nil {
					d.Set("partition_key_version", version)
				}
			}

			if ukp := res.UniqueKeyPolicy; ukp != nil {
				if err := d.Set("unique_key", flattenCosmosSQLContainerUniqueKeys(ukp.UniqueKeys)); err != nil {
					return fmt.Errorf("Error setting `unique_key`: %+v", err)
				}
			}

			if defaultTTL := res.DefaultTTL; defaultTTL != nil {
				d.Set("default_ttl", defaultTTL)
			}

			if indexingPolicy := res.IndexingPolicy; indexingPolicy != nil {
				d.Set("indexing_policy", common.FlattenAzureRmCosmosDbIndexingPolicy(indexingPolicy))
			}
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
			throughputResp, err := client.GetSQLContainerThroughput(ctx, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName)
			if err != nil {
				if !utils.ResponseWasNotFound(throughputResp.Response) {
					return fmt.Errorf("Error reading Throughput on Cosmos SQL Container %s (Account: %q, Database: %q) ID: %v", id.ContainerName, id.DatabaseAccountName, id.SqlDatabaseName, err)
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

func resourceArmCosmosDbSQLContainerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlContainerID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.DeleteSQLContainer(ctx, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting Cosmos SQL Container %q (Account: %q): %+v", id.SqlDatabaseName, id.ContainerName, err)
		}
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting on delete future for Cosmos SQL Container %q (Account: %q): %+v", id.SqlDatabaseName, id.DatabaseAccountName, err)
	}

	return nil
}

func expandCosmosSQLContainerUniqueKeys(s *schema.Set) *[]documentdb.UniqueKey {
	i := s.List()
	if len(i) == 0 || i[0] == nil {
		return nil
	}

	keys := make([]documentdb.UniqueKey, 0)
	for _, k := range i {
		key := k.(map[string]interface{})

		paths := key["paths"].(*schema.Set).List()
		if len(paths) == 0 {
			continue
		}

		keys = append(keys, documentdb.UniqueKey{
			Paths: utils.ExpandStringSlice(paths),
		})
	}

	return &keys
}

func flattenCosmosSQLContainerUniqueKeys(keys *[]documentdb.UniqueKey) *[]map[string]interface{} {
	if keys == nil {
		return nil
	}

	slice := make([]map[string]interface{}, 0)
	for _, k := range *keys {
		if k.Paths == nil {
			continue
		}

		slice = append(slice, map[string]interface{}{
			"paths": *k.Paths,
		})
	}

	return &slice
}
