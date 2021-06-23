package cosmos

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-01-15/documentdb"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/migration"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceCosmosDbSQLContainer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCosmosDbSQLContainerCreate,
		Read:   resourceCosmosDbSQLContainerRead,
		Update: resourceCosmosDbSQLContainerUpdate,
		Delete: resourceCosmosDbSQLContainerDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SqlContainerV0ToV1{},
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

			"resource_group_name": azure.SchemaResourceGroupName(),

			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosAccountName,
			},

			"database_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosEntityName,
			},

			"partition_key_path": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"partition_key_version": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 2),
			},

			"conflict_resolution_policy": common.ConflictResolutionPolicy(),

			"throughput": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.CosmosThroughput,
			},

			"autoscale_settings": common.DatabaseAutoscaleSettingsSchema(),

			"analytical_storage_ttl": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(-1),
			},

			"default_ttl": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(-1),
			},

			"unique_key": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"paths": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							ForceNew: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
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

func resourceCosmosDbSQLContainerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	database := d.Get("database_name").(string)
	account := d.Get("account_name").(string)
	partitionkeypaths := d.Get("partition_key_path").(string)

	existing, err := client.GetSQLContainer(ctx, resourceGroup, account, database, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of creating Cosmos SQL Container %q (Account: %q, Database: %q): %+v", name, account, database, err)
		}
	} else {
		if existing.ID == nil && *existing.ID == "" {
			return fmt.Errorf("generating import ID for Cosmos SQL Container %q (Account: %q, Database: %q)", name, account, database)
		}

		return tf.ImportAsExistsError("azurerm_cosmosdb_sql_container", *existing.ID)
	}

	indexingPolicy := common.ExpandAzureRmCosmosDbIndexingPolicy(d)
	err = common.ValidateAzureRmCosmosDbIndexingPolicy(indexingPolicy)
	if err != nil {
		return fmt.Errorf("generating indexing policy for Cosmos SQL Container %q (Account: %q, Database: %q)", name, account, database)
	}

	db := documentdb.SQLContainerCreateUpdateParameters{
		SQLContainerCreateUpdateProperties: &documentdb.SQLContainerCreateUpdateProperties{
			Resource: &documentdb.SQLContainerResource{
				ID:                       &name,
				IndexingPolicy:           indexingPolicy,
				ConflictResolutionPolicy: common.ExpandCosmosDbConflicResolutionPolicy(d.Get("conflict_resolution_policy").([]interface{})),
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	if partitionkeypaths != "" {
		db.SQLContainerCreateUpdateProperties.Resource.PartitionKey = &documentdb.ContainerPartitionKey{
			Paths: &[]string{partitionkeypaths},
			Kind:  documentdb.PartitionKindHash,
		}

		if partitionKeyVersion, ok := d.GetOk("partition_key_version"); ok {
			db.SQLContainerCreateUpdateProperties.Resource.PartitionKey.Version = utils.Int32(int32(partitionKeyVersion.(int)))
		}
	}

	if keys := expandCosmosSQLContainerUniqueKeys(d.Get("unique_key").(*pluginsdk.Set)); keys != nil {
		db.SQLContainerCreateUpdateProperties.Resource.UniqueKeyPolicy = &documentdb.UniqueKeyPolicy{
			UniqueKeys: keys,
		}
	}

	if analyticalStorageTTL, ok := d.GetOk("analytical_storage_ttl"); ok {
		db.SQLContainerCreateUpdateProperties.Resource.AnalyticalStorageTTL = utils.Int64(int64(analyticalStorageTTL.(int)))
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
		return fmt.Errorf("issuing create/update request for Cosmos SQL Container %q (Account: %q, Database: %q): %+v", name, account, database, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on create/update future for Cosmos SQL Container %q (Account: %q, Database: %q): %+v", name, account, database, err)
	}

	resp, err := client.GetSQLContainer(ctx, resourceGroup, account, database, name)
	if err != nil {
		return fmt.Errorf("making get request for Cosmos SQL Container %q (Account: %q, Database: %q): %+v", name, account, database, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("getting ID from Cosmos SQL Container %q (Account: %q, Database: %q)", name, account, database)
	}

	d.SetId(*resp.ID)

	return resourceCosmosDbSQLContainerRead(d, meta)
}

func resourceCosmosDbSQLContainerUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlContainerID(d.Id())
	if err != nil {
		return err
	}

	err = common.CheckForChangeFromAutoscaleAndManualThroughput(d)
	if err != nil {
		return fmt.Errorf("updating Cosmos SQL Container %q (Account: %q, Database: %q): %+v", id.ContainerName, id.DatabaseAccountName, id.SqlDatabaseName, err)
	}

	partitionkeypaths := d.Get("partition_key_path").(string)

	indexingPolicy := common.ExpandAzureRmCosmosDbIndexingPolicy(d)
	err = common.ValidateAzureRmCosmosDbIndexingPolicy(indexingPolicy)
	if err != nil {
		return fmt.Errorf("updating Cosmos SQL Container %q (Account: %q, Database: %q): %+v", id.ContainerName, id.DatabaseAccountName, id.SqlDatabaseName, err)
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
		db.SQLContainerCreateUpdateProperties.Resource.PartitionKey = &documentdb.ContainerPartitionKey{
			Paths: &[]string{partitionkeypaths},
			Kind:  documentdb.PartitionKindHash,
		}

		if partitionKeyVersion, ok := d.GetOk("partition_key_version"); ok {
			db.SQLContainerCreateUpdateProperties.Resource.PartitionKey.Version = utils.Int32(int32(partitionKeyVersion.(int)))
		}
	}

	if keys := expandCosmosSQLContainerUniqueKeys(d.Get("unique_key").(*pluginsdk.Set)); keys != nil {
		db.SQLContainerCreateUpdateProperties.Resource.UniqueKeyPolicy = &documentdb.UniqueKeyPolicy{
			UniqueKeys: keys,
		}
	}

	if analyticalStorageTTL, ok := d.GetOk("analytical_storage_ttl"); ok {
		db.SQLContainerCreateUpdateProperties.Resource.AnalyticalStorageTTL = utils.Int64(int64(analyticalStorageTTL.(int)))
	}

	if defaultTTL, hasTTL := d.GetOk("default_ttl"); hasTTL {
		db.SQLContainerCreateUpdateProperties.Resource.DefaultTTL = utils.Int32(int32(defaultTTL.(int)))
	}

	future, err := client.CreateUpdateSQLContainer(ctx, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName, db)
	if err != nil {
		return fmt.Errorf("issuing create/update request for Cosmos SQL Container %q (Account: %q, Database: %q): %+v", id.ContainerName, id.DatabaseAccountName, id.SqlDatabaseName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on create/update future for Cosmos SQL Container %q (Account: %q, Database: %q): %+v", id.ContainerName, id.DatabaseAccountName, id.SqlDatabaseName, err)
	}

	if common.HasThroughputChange(d) {
		throughputParameters := common.ExpandCosmosDBThroughputSettingsUpdateParameters(d)
		throughputFuture, err := client.UpdateSQLContainerThroughput(ctx, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName, *throughputParameters)
		if err != nil {
			if response.WasNotFound(throughputFuture.Response()) {
				return fmt.Errorf("setting Throughput for Cosmos SQL Container %q (Account: %q, Database: %q): %+v - "+
					"If the collection has not been created with an initial throughput, you cannot configure it later", id.ContainerName, id.DatabaseAccountName, id.SqlDatabaseName, err)
			}
		}

		if err = throughputFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting on ThroughputUpdate future for Cosmos Container %q (Account: %q, Database: %q): %+v", id.ContainerName, id.DatabaseAccountName, id.SqlDatabaseName, err)
		}
	}

	return resourceCosmosDbSQLContainerRead(d, meta)
}

func resourceCosmosDbSQLContainerRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

		return fmt.Errorf("reading Cosmos SQL Container %q (Account: %q): %+v", id.SqlDatabaseName, id.ContainerName, err)
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
						return fmt.Errorf("reading PartitionKey Paths, more then 1 returned")
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
					return fmt.Errorf("setting `unique_key`: %+v", err)
				}
			}

			if analyticalStorageTTL := res.AnalyticalStorageTTL; analyticalStorageTTL != nil {
				d.Set("analytical_storage_ttl", analyticalStorageTTL)
			}

			if defaultTTL := res.DefaultTTL; defaultTTL != nil {
				d.Set("default_ttl", defaultTTL)
			}

			if indexingPolicy := res.IndexingPolicy; indexingPolicy != nil {
				d.Set("indexing_policy", common.FlattenAzureRmCosmosDbIndexingPolicy(indexingPolicy))
			}

			if err := d.Set("conflict_resolution_policy", common.FlattenCosmosDbConflictResolutionPolicy(res.ConflictResolutionPolicy)); err != nil {
				return fmt.Errorf("setting `conflict_resolution_policy`: %+v", err)
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

	// if the cosmos account is serverless calling the get throughput api would yield an error
	if !isServerlessCapacityMode(accResp) {
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

	return nil
}

func resourceCosmosDbSQLContainerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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

func expandCosmosSQLContainerUniqueKeys(s *pluginsdk.Set) *[]documentdb.UniqueKey {
	i := s.List()
	if len(i) == 0 || i[0] == nil {
		return nil
	}

	keys := make([]documentdb.UniqueKey, 0)
	for _, k := range i {
		key := k.(map[string]interface{})

		paths := key["paths"].(*pluginsdk.Set).List()
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
