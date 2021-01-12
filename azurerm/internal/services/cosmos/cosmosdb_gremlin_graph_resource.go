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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/migration"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceCosmosDbGremlinGraph() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmosDbGremlinGraphCreate,
		Read:   resourceCosmosDbGremlinGraphRead,
		Update: resourceCosmosDbGremlinGraphUpdate,
		Delete: resourceCosmosDbGremlinGraphDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    migration.ResourceGremlinGraphUpgradeV0Schema().CoreConfigSchema().ImpliedType(),
				Upgrade: migration.ResourceGremlinGraphStateUpgradeV0ToV1,
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

			"default_ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"throughput": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.CosmosThroughput,
			},

			"autoscale_settings": common.ContainerAutoscaleSettingsSchema(),

			"partition_key_path": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"index_policy": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"automatic": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"indexing_mode": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference, // Open issue https://github.com/Azure/azure-sdk-for-go/issues/6603
							ValidateFunc: validation.StringInSlice([]string{
								string(documentdb.Consistent),
								string(documentdb.Lazy),
								string(documentdb.None),
							}, false),
						},

						"included_paths": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							Set: schema.HashString,
						},

						"excluded_paths": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							Set: schema.HashString,
						},
					},
				},
			},

			"conflict_resolution_policy": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(documentdb.LastWriterWins),
								string(documentdb.Custom),
							}, false),
						},

						"conflict_resolution_path": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"conflict_resolution_procedure": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
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
		},
	}
}

func resourceCosmosDbGremlinGraphCreate(d *schema.ResourceData, meta interface{}) error {
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
				ConflictResolutionPolicy: expandAzureRmCosmosDbGremlinGraphConflicResolutionPolicy(d),
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	if partitionkeypaths != "" {
		db.GremlinGraphCreateUpdateProperties.Resource.PartitionKey = &documentdb.ContainerPartitionKey{
			Paths: &[]string{partitionkeypaths},
		}
	}

	if keys := expandAzureRmCosmosDbGremlinGraphUniqueKeys(d.Get("unique_key").(*schema.Set)); keys != nil {
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

func resourceCosmosDbGremlinGraphUpdate(d *schema.ResourceData, meta interface{}) error {
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
				ID:                       &id.GraphName,
				IndexingPolicy:           expandAzureRmCosmosDbGrelinGraphIndexingPolicy(d),
				ConflictResolutionPolicy: expandAzureRmCosmosDbGremlinGraphConflicResolutionPolicy(d),
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	if partitionkeypaths != "" {
		db.GremlinGraphCreateUpdateProperties.Resource.PartitionKey = &documentdb.ContainerPartitionKey{
			Paths: &[]string{partitionkeypaths},
		}
	}

	if keys := expandAzureRmCosmosDbGremlinGraphUniqueKeys(d.Get("unique_key").(*schema.Set)); keys != nil {
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

func resourceCosmosDbGremlinGraphRead(d *schema.ResourceData, meta interface{}) error {
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
			}

			if ip := props.IndexingPolicy; ip != nil {
				if err := d.Set("index_policy", flattenAzureRmCosmosDBGremlinGraphIndexingPolicy(props.IndexingPolicy)); err != nil {
					return fmt.Errorf("Error setting `index_policy`: %+v", err)
				}
			}

			if crp := props.ConflictResolutionPolicy; crp != nil {
				if err := d.Set("conflict_resolution_policy", flattenAzureRmCosmosDbGremlinGraphConflictResolutionPolicy(props.ConflictResolutionPolicy)); err != nil {
					return fmt.Errorf("Error setting `conflict_resolution_policy`: %+v", err)
				}
			}

			if ukp := props.UniqueKeyPolicy; ukp != nil {
				if err := d.Set("unique_key", flattenCosmosGremlinGraphUniqueKeys(ukp.UniqueKeys)); err != nil {
					return fmt.Errorf("Error setting `unique_key`: %+v", err)
				}
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

func resourceCosmosDbGremlinGraphDelete(d *schema.ResourceData, meta interface{}) error {
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

func expandAzureRmCosmosDbGrelinGraphIndexingPolicy(d *schema.ResourceData) *documentdb.IndexingPolicy {
	i := d.Get("index_policy").([]interface{})
	if len(i) == 0 || i[0] == nil {
		return nil
	}

	input := i[0].(map[string]interface{})
	indexingPolicy := input["indexing_mode"].(string)
	policy := &documentdb.IndexingPolicy{
		IndexingMode:  documentdb.IndexingMode(indexingPolicy),
		IncludedPaths: expandAzureRmCosmosDbGrelimGraphIncludedPath(input),
		ExcludedPaths: expandAzureRmCosmosDbGremlinGraphExcludedPath(input),
	}

	if automatic, ok := input["automatic"].(bool); ok {
		policy.Automatic = utils.Bool(automatic)
	}

	return policy
}

func expandAzureRmCosmosDbGremlinGraphConflicResolutionPolicy(d *schema.ResourceData) *documentdb.ConflictResolutionPolicy {
	i := d.Get("conflict_resolution_policy").([]interface{})
	if len(i) == 0 || i[0] == nil {
		return nil
	}

	input := i[0].(map[string]interface{})
	conflictResolutionMode := input["mode"].(string)
	conflict := &documentdb.ConflictResolutionPolicy{
		Mode: documentdb.ConflictResolutionMode(conflictResolutionMode),
	}

	if conflictResolutionPath, ok := input["conflict_resolution_path"].(string); ok {
		conflict.ConflictResolutionPath = utils.String(conflictResolutionPath)
	}

	if conflictResolutionProcedure, ok := input["conflict_resolution_procedure"].(string); ok {
		conflict.ConflictResolutionProcedure = utils.String(conflictResolutionProcedure)
	}

	return conflict
}

func expandAzureRmCosmosDbGrelimGraphIncludedPath(input map[string]interface{}) *[]documentdb.IncludedPath {
	includedPath := input["included_paths"].(*schema.Set).List()
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
	excludedPath := input["excluded_paths"].(*schema.Set).List()
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

func expandAzureRmCosmosDbGremlinGraphUniqueKeys(s *schema.Set) *[]documentdb.UniqueKey {
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

func flattenAzureRmCosmosDBGremlinGraphIndexingPolicy(input *documentdb.IndexingPolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	indexPolicy := make(map[string]interface{})

	indexPolicy["automatic"] = input.Automatic
	indexPolicy["indexing_mode"] = string(input.IndexingMode)
	indexPolicy["included_paths"] = schema.NewSet(schema.HashString, flattenAzureRmCosmosDBGremlinGraphIncludedPaths(input.IncludedPaths))
	indexPolicy["excluded_paths"] = schema.NewSet(schema.HashString, flattenAzureRmCosmosDBGremlinGraphExcludedPaths(input.ExcludedPaths))

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

func flattenAzureRmCosmosDbGremlinGraphConflictResolutionPolicy(input *documentdb.ConflictResolutionPolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	conflictResolutionPolicy := make(map[string]interface{})

	conflictResolutionPolicy["mode"] = string(input.Mode)
	conflictResolutionPolicy["conflict_resolution_path"] = input.ConflictResolutionPath
	conflictResolutionPolicy["conflict_resolution_procedure"] = input.ConflictResolutionProcedure

	return []interface{}{conflictResolutionPolicy}
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
