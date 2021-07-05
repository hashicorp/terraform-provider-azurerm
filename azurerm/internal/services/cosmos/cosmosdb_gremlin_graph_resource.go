package cosmos

import (
	"fmt"
	"log"
	"strings"
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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceCosmosDbGremlinGraph() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCosmosDbGremlinGraphCreate,
		Read:   resourceCosmosDbGremlinGraphRead,
		Update: resourceCosmosDbGremlinGraphUpdate,
		Delete: resourceCosmosDbGremlinGraphDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.GremlinGraphV0ToV1{},
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

			"default_ttl": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Computed: true,
			},

			"throughput": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.CosmosThroughput,
			},

			"autoscale_settings": common.DatabaseAutoscaleSettingsSchema(),

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

			"index_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"automatic": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						// case change in 2021-01-15, issue https://github.com/Azure/azure-rest-api-specs/issues/14051
						// todo: change to SDK constants and remove translation code in 3.0
						"indexing_mode": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference, // Open issue https://github.com/Azure/azure-sdk-for-go/issues/6603
							ValidateFunc: validation.StringInSlice([]string{
								"Consistent",
								"Lazy",
								"None",
							}, false),
						},

						"included_paths": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							Set: pluginsdk.HashString,
						},

						"excluded_paths": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							Set: pluginsdk.HashString,
						},

						"composite_index": common.CosmosDbIndexingPolicyCompositeIndexSchema(),

						"spatial_index": common.CosmosDbIndexingPolicySpatialIndexSchema(),
					},
				},
			},

			"conflict_resolution_policy": common.ConflictResolutionPolicy(),

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
		},
	}
}

func resourceCosmosDbGremlinGraphCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.GremlinClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	database := d.Get("database_name").(string)
	account := d.Get("account_name").(string)
	partitionkeypaths := d.Get("partition_key_path").(string)

	existing, err := client.GetGremlinGraph(ctx, resourceGroup, account, database, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of creating Cosmos Gremlin Graph %q (Account: %q, Database: %q): %+v", name, account, database, err)
		}
	} else {
		if existing.ID == nil && *existing.ID == "" {
			return fmt.Errorf("Error generating import ID for Cosmos Gremlin Graph %q (Account: %q, Database: %q)", name, account, database)
		}

		return tf.ImportAsExistsError("azurerm_cosmosdb_gremlin_graph", *existing.ID)
	}

	db := documentdb.GremlinGraphCreateUpdateParameters{
		GremlinGraphCreateUpdateProperties: &documentdb.GremlinGraphCreateUpdateProperties{
			Resource: &documentdb.GremlinGraphResource{
				ID:                       &name,
				IndexingPolicy:           expandAzureRmCosmosDbGrelinGraphIndexingPolicy(d),
				ConflictResolutionPolicy: common.ExpandCosmosDbConflicResolutionPolicy(d.Get("conflict_resolution_policy").([]interface{})),
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	if partitionkeypaths != "" {
		db.GremlinGraphCreateUpdateProperties.Resource.PartitionKey = &documentdb.ContainerPartitionKey{
			Paths: &[]string{partitionkeypaths},
			Kind:  documentdb.PartitionKindHash,
		}
		if partitionKeyVersion, ok := d.GetOk("partition_key_version"); ok {
			db.GremlinGraphCreateUpdateProperties.Resource.PartitionKey.Version = utils.Int32(int32(partitionKeyVersion.(int)))
		}
	}

	if keys := expandAzureRmCosmosDbGremlinGraphUniqueKeys(d.Get("unique_key").(*pluginsdk.Set)); keys != nil {
		db.GremlinGraphCreateUpdateProperties.Resource.UniqueKeyPolicy = &documentdb.UniqueKeyPolicy{
			UniqueKeys: keys,
		}
	}

	if defaultTTL, hasDefaultTTL := d.GetOk("default_ttl"); hasDefaultTTL {
		if defaultTTL != 0 {
			db.GremlinGraphCreateUpdateProperties.Resource.DefaultTTL = utils.Int32(int32(defaultTTL.(int)))
		}
	}

	if throughput, hasThroughput := d.GetOk("throughput"); hasThroughput {
		if throughput != 0 {
			db.GremlinGraphCreateUpdateProperties.Options.Throughput = common.ConvertThroughputFromResourceData(throughput)
		}
	}

	if _, hasAutoscaleSettings := d.GetOk("autoscale_settings"); hasAutoscaleSettings {
		db.GremlinGraphCreateUpdateProperties.Options.AutoscaleSettings = common.ExpandCosmosDbAutoscaleSettings(d)
	}

	future, err := client.CreateUpdateGremlinGraph(ctx, resourceGroup, account, database, name, db)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Cosmos Gremlin Graph %q (Account: %q, Database: %q): %+v", name, account, database, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Cosmos Gremlin Graph%q (Account: %q, Database: %q): %+v", name, account, database, err)
	}

	resp, err := client.GetGremlinGraph(ctx, resourceGroup, account, database, name)
	if err != nil {
		return fmt.Errorf("Error making get request for Cosmos Gremlin Graph %q (Account: %q, Database: %q): %+v", name, account, database, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Error getting ID from Cosmos Gremlin Graph %q (Account: %q, Database: %q)", name, account, database)
	}

	d.SetId(*resp.ID)

	return resourceCosmosDbGremlinGraphRead(d, meta)
}

func resourceCosmosDbGremlinGraphUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.GremlinClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.GremlinGraphID(d.Id())
	if err != nil {
		return err
	}

	err = common.CheckForChangeFromAutoscaleAndManualThroughput(d)
	if err != nil {
		return fmt.Errorf("Error updating Cosmos Gremlin Graph %q (Account: %q, Database: %q): %+v", id.GraphName, id.DatabaseAccountName, id.GremlinDatabaseName, err)
	}

	partitionkeypaths := d.Get("partition_key_path").(string)

	db := documentdb.GremlinGraphCreateUpdateParameters{
		GremlinGraphCreateUpdateProperties: &documentdb.GremlinGraphCreateUpdateProperties{
			Resource: &documentdb.GremlinGraphResource{
				ID:             &id.GraphName,
				IndexingPolicy: expandAzureRmCosmosDbGrelinGraphIndexingPolicy(d),
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	if partitionkeypaths != "" {
		db.GremlinGraphCreateUpdateProperties.Resource.PartitionKey = &documentdb.ContainerPartitionKey{
			Paths: &[]string{partitionkeypaths},
			Kind:  documentdb.PartitionKindHash,
		}

		if partitionKeyVersion, ok := d.GetOk("partition_key_version"); ok {
			db.GremlinGraphCreateUpdateProperties.Resource.PartitionKey.Version = utils.Int32(int32(partitionKeyVersion.(int)))
		}
	}

	if keys := expandAzureRmCosmosDbGremlinGraphUniqueKeys(d.Get("unique_key").(*pluginsdk.Set)); keys != nil {
		db.GremlinGraphCreateUpdateProperties.Resource.UniqueKeyPolicy = &documentdb.UniqueKeyPolicy{
			UniqueKeys: keys,
		}
	}

	if defaultTTL, hasDefaultTTL := d.GetOk("default_ttl"); hasDefaultTTL {
		db.GremlinGraphCreateUpdateProperties.Resource.DefaultTTL = utils.Int32(int32(defaultTTL.(int)))
	}

	future, err := client.CreateUpdateGremlinGraph(ctx, id.ResourceGroup, id.DatabaseAccountName, id.GremlinDatabaseName, id.GraphName, db)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Cosmos Gremlin Graph %q (Account: %q, Database: %q): %+v", id.GraphName, id.DatabaseAccountName, id.GremlinDatabaseName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Cosmos Gremlin Graph %q (Account: %q, Database: %q): %+v", id.GraphName, id.DatabaseAccountName, id.GremlinDatabaseName, err)
	}

	if common.HasThroughputChange(d) {
		throughputParameters := common.ExpandCosmosDBThroughputSettingsUpdateParameters(d)
		throughputFuture, err := client.UpdateGremlinGraphThroughput(ctx, id.ResourceGroup, id.DatabaseAccountName, id.GremlinDatabaseName, id.GraphName, *throughputParameters)
		if err != nil {
			if response.WasNotFound(throughputFuture.Response()) {
				return fmt.Errorf("Error setting Throughput for Cosmos Gremlin Graph %q (Account: %q, Database: %q): %+v - "+
					"If the graph has not been created with an initial throughput, you cannot configure it later.", id.GraphName, id.DatabaseAccountName, id.GremlinDatabaseName, err)
			}
		}

		if err = throughputFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting on ThroughputUpdate future for Cosmos Gremlin Graph %q (Account: %q, Database: %q): %+v", id.GraphName, id.DatabaseAccountName, id.GremlinDatabaseName, err)
		}
	}

	return resourceCosmosDbGremlinGraphRead(d, meta)
}

func resourceCosmosDbGremlinGraphRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.GremlinClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.GremlinGraphID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetGremlinGraph(ctx, id.ResourceGroup, id.DatabaseAccountName, id.GremlinDatabaseName, id.GraphName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Cosmos Gremlin Graph %q (Account: %q) - removing from state", id.GraphName, id.DatabaseAccountName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Cosmos Gremlin Graph %q (Account: %q): %+v", id.GraphName, id.DatabaseAccountName, err)
	}

	d.Set("name", id.GraphName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_name", id.DatabaseAccountName)
	d.Set("database_name", id.GremlinDatabaseName)

	if graphProperties := resp.GremlinGraphGetProperties; graphProperties != nil {
		if props := graphProperties.Resource; props != nil {
			if pk := props.PartitionKey; pk != nil {
				if paths := pk.Paths; paths != nil {
					if len(*paths) > 1 {
						return fmt.Errorf("Error reading PartitionKey Paths, more than 1 returned")
					} else if len(*paths) == 1 {
						d.Set("partition_key_path", (*paths)[0])
					}
				}

				if version := pk.Version; version != nil {
					d.Set("partition_key_version", version)
				}
			}

			if ip := props.IndexingPolicy; ip != nil {
				if err := d.Set("index_policy", flattenAzureRmCosmosDBGremlinGraphIndexingPolicy(props.IndexingPolicy)); err != nil {
					return fmt.Errorf("Error setting `index_policy`: %+v", err)
				}
			}

			if crp := props.ConflictResolutionPolicy; crp != nil {
				if err := d.Set("conflict_resolution_policy", common.FlattenCosmosDbConflictResolutionPolicy(crp)); err != nil {
					return fmt.Errorf("Error setting `conflict_resolution_policy`: %+v", err)
				}
			}

			if ukp := props.UniqueKeyPolicy; ukp != nil {
				if err := d.Set("unique_key", flattenCosmosGremlinGraphUniqueKeys(ukp.UniqueKeys)); err != nil {
					return fmt.Errorf("Error setting `unique_key`: %+v", err)
				}
			}

			if defaultTTL := props.DefaultTTL; defaultTTL != nil {
				d.Set("default_ttl", defaultTTL)
			}
		}
	}

	throughputResp, err := client.GetGremlinGraphThroughput(ctx, id.ResourceGroup, id.DatabaseAccountName, id.GremlinDatabaseName, id.GraphName)
	if err != nil {
		if !utils.ResponseWasNotFound(throughputResp.Response) {
			return fmt.Errorf("Error reading Throughput on Gremlin Graph %q (Account: %q, Database: %q) ID: %v", id.GraphName, id.DatabaseAccountName, id.GremlinDatabaseName, err)
		} else {
			d.Set("throughput", nil)
			d.Set("autoscale_settings", nil)
		}
	} else {
		common.SetResourceDataThroughputFromResponse(throughputResp, d)
	}

	return nil
}

func resourceCosmosDbGremlinGraphDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.GremlinClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.GremlinGraphID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.DeleteGremlinGraph(ctx, id.ResourceGroup, id.DatabaseAccountName, id.GremlinDatabaseName, id.GraphName)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting Cosmos Gremlin Graph %q (Account: %q): %+v", id.GremlinDatabaseName, id.GraphName, err)
		}
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on delete future for Comos Gremlin Graph %q (Account: %q): %+v", id.GremlinDatabaseName, id.DatabaseAccountName, err)
	}

	return nil
}

func expandAzureRmCosmosDbGrelinGraphIndexingPolicy(d *pluginsdk.ResourceData) *documentdb.IndexingPolicy {
	i := d.Get("index_policy").([]interface{})
	if len(i) == 0 || i[0] == nil {
		return nil
	}

	input := i[0].(map[string]interface{})
	indexingPolicy := input["indexing_mode"].(string)
	policy := &documentdb.IndexingPolicy{
		IndexingMode:  documentdb.IndexingMode(strings.ToLower(indexingPolicy)),
		IncludedPaths: expandAzureRmCosmosDbGrelimGraphIncludedPath(input),
		ExcludedPaths: expandAzureRmCosmosDbGremlinGraphExcludedPath(input),
	}
	if v, ok := input["composite_index"].([]interface{}); ok {
		policy.CompositeIndexes = common.ExpandAzureRmCosmosDBIndexingPolicyCompositeIndexes(v)
	}

	policy.SpatialIndexes = common.ExpandAzureRmCosmosDBIndexingPolicySpatialIndexes(input["spatial_index"].([]interface{}))

	if automatic, ok := input["automatic"].(bool); ok {
		policy.Automatic = utils.Bool(automatic)
	}

	return policy
}

func expandAzureRmCosmosDbGrelimGraphIncludedPath(input map[string]interface{}) *[]documentdb.IncludedPath {
	includedPath := input["included_paths"].(*pluginsdk.Set).List()
	paths := make([]documentdb.IncludedPath, len(includedPath))

	for i, pathConfig := range includedPath {
		attrs := pathConfig.(string)
		path := documentdb.IncludedPath{
			Path: utils.String(attrs),
		}
		paths[i] = path
	}

	return &paths
}

func expandAzureRmCosmosDbGremlinGraphExcludedPath(input map[string]interface{}) *[]documentdb.ExcludedPath {
	excludedPath := input["excluded_paths"].(*pluginsdk.Set).List()
	paths := make([]documentdb.ExcludedPath, len(excludedPath))

	for i, pathConfig := range excludedPath {
		attrs := pathConfig.(string)
		path := documentdb.ExcludedPath{
			Path: utils.String(attrs),
		}
		paths[i] = path
	}

	return &paths
}

func expandAzureRmCosmosDbGremlinGraphUniqueKeys(s *pluginsdk.Set) *[]documentdb.UniqueKey {
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

func flattenAzureRmCosmosDBGremlinGraphIndexingPolicy(input *documentdb.IndexingPolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	indexPolicy := make(map[string]interface{})

	indexPolicy["automatic"] = input.Automatic
	indexPolicy["indexing_mode"] = strings.Title(string(input.IndexingMode))
	indexPolicy["included_paths"] = pluginsdk.NewSet(pluginsdk.HashString, flattenAzureRmCosmosDBGremlinGraphIncludedPaths(input.IncludedPaths))
	indexPolicy["excluded_paths"] = pluginsdk.NewSet(pluginsdk.HashString, flattenAzureRmCosmosDBGremlinGraphExcludedPaths(input.ExcludedPaths))
	indexPolicy["composite_index"] = common.FlattenCosmosDBIndexingPolicyCompositeIndexes(input.CompositeIndexes)
	indexPolicy["spatial_index"] = common.FlattenCosmosDBIndexingPolicySpatialIndexes(input.SpatialIndexes)

	return []interface{}{indexPolicy}
}

func flattenAzureRmCosmosDBGremlinGraphIncludedPaths(input *[]documentdb.IncludedPath) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	includedPaths := make([]interface{}, 0)
	for _, includedPath := range *input {
		if includedPath.Path == nil {
			continue
		}

		includedPaths = append(includedPaths, *includedPath.Path)
	}

	return includedPaths
}

func flattenAzureRmCosmosDBGremlinGraphExcludedPaths(input *[]documentdb.ExcludedPath) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	excludedPaths := make([]interface{}, 0)
	for _, excludedPath := range *input {
		if excludedPath.Path == nil {
			continue
		}

		excludedPaths = append(excludedPaths, *excludedPath.Path)
	}

	return excludedPaths
}

func flattenCosmosGremlinGraphUniqueKeys(keys *[]documentdb.UniqueKey) *[]map[string]interface{} {
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
