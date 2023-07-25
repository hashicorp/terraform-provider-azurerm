// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2023-04-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCosmosDbGremlinGraph() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCosmosDbGremlinGraphCreate,
		Read:   resourceCosmosDbGremlinGraphRead,
		Update: resourceCosmosDbGremlinGraphUpdate,
		Delete: resourceCosmosDbGremlinGraphDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.GremlinGraphID(id)
			return err
		}),

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

			"resource_group_name": commonschema.ResourceGroupName(),

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

			"analytical_storage_ttl": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				ValidateFunc: validation.All(
					validation.IntBetween(-1, 2147483647),
					validation.IntNotInSlice([]int{0}),
				),
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
						"indexing_mode": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(cosmosdb.IndexingModeConsistent),
								string(cosmosdb.IndexingModeNone),
								string(cosmosdb.IndexingModeLazy),
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

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			// `analytical_storage_ttl` can't be disabled once it's enabled
			pluginsdk.ForceNewIfChange("analytical_storage_ttl", func(ctx context.Context, old, new, _ interface{}) bool {
				return (old.(int) == -1 || (old.(int) >= 1 && old.(int) <= 2147483647)) && new.(int) == 0
			}),
		),
	}
}

func resourceCosmosDbGremlinGraphCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := cosmosdb.NewGraphID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("database_name").(string), d.Get("name").(string))
	partitionkeypaths := d.Get("partition_key_path").(string)

	existing, err := client.GremlinResourcesGetGremlinGraph(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_cosmosdb_gremlin_graph", id.ID())
	}

	db := cosmosdb.GremlinGraphCreateUpdateParameters{
		Properties: cosmosdb.GremlinGraphCreateUpdateProperties{
			Resource: cosmosdb.GremlinGraphResource{
				Id:                       id.GraphName,
				IndexingPolicy:           expandAzureRmCosmosDbGrelinGraphIndexingPolicy(d),
				ConflictResolutionPolicy: common.ExpandCosmosDbConflicResolutionPolicy(d.Get("conflict_resolution_policy").([]interface{})),
			},
			Options: &cosmosdb.CreateUpdateOptions{},
		},
	}

	if v, ok := d.GetOk("analytical_storage_ttl"); ok {
		db.Properties.Resource.AnalyticalStorageTtl = utils.Int64(int64(v.(int)))
	}

	if partitionkeypaths != "" {
		partitionKindHash := cosmosdb.PartitionKindHash
		db.Properties.Resource.PartitionKey = &cosmosdb.ContainerPartitionKey{
			Paths: &[]string{partitionkeypaths},
			Kind:  &partitionKindHash,
		}
		if partitionKeyVersion, ok := d.GetOk("partition_key_version"); ok {
			db.Properties.Resource.PartitionKey.Version = utils.Int64(int64(partitionKeyVersion.(int)))
		}
	}

	if keys := expandAzureRmCosmosDbGremlinGraphUniqueKeys(d.Get("unique_key").(*pluginsdk.Set)); keys != nil {
		db.Properties.Resource.UniqueKeyPolicy = &cosmosdb.UniqueKeyPolicy{
			UniqueKeys: keys,
		}
	}

	if defaultTTL, hasDefaultTTL := d.GetOk("default_ttl"); hasDefaultTTL {
		if defaultTTL != 0 {
			db.Properties.Resource.DefaultTtl = utils.Int64(int64(defaultTTL.(int)))
		}
	}

	if throughput, hasThroughput := d.GetOk("throughput"); hasThroughput {
		if throughput != 0 {
			db.Properties.Options.Throughput = common.ConvertThroughputFromResourceData(throughput)
		}
	}

	if _, hasAutoscaleSettings := d.GetOk("autoscale_settings"); hasAutoscaleSettings {
		db.Properties.Options.AutoScaleSettings = common.ExpandCosmosDbAutoscaleSettings(d)
	}

	err = client.GremlinResourcesCreateUpdateGremlinGraphThenPoll(ctx, id, db)
	if err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosDbGremlinGraphRead(d, meta)
}

func resourceCosmosDbGremlinGraphUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseGraphID(d.Id())
	if err != nil {
		return err
	}

	err = common.CheckForChangeFromAutoscaleAndManualThroughput(d)
	if err != nil {
		return fmt.Errorf("updating Cosmos Gremlin Graph %q (Account: %q, Database: %q): %+v", id.GraphName, id.DatabaseAccountName, id.GremlinDatabaseName, err)
	}

	partitionkeypaths := d.Get("partition_key_path").(string)

	db := cosmosdb.GremlinGraphCreateUpdateParameters{
		Properties: cosmosdb.GremlinGraphCreateUpdateProperties{
			Resource: cosmosdb.GremlinGraphResource{
				Id:             id.GraphName,
				IndexingPolicy: expandAzureRmCosmosDbGrelinGraphIndexingPolicy(d),
			},
			Options: &cosmosdb.CreateUpdateOptions{},
		},
	}

	if partitionkeypaths != "" {
		partitionKindHash := cosmosdb.PartitionKindHash
		db.Properties.Resource.PartitionKey = &cosmosdb.ContainerPartitionKey{
			Paths: &[]string{partitionkeypaths},
			Kind:  &partitionKindHash,
		}

		if partitionKeyVersion, ok := d.GetOk("partition_key_version"); ok {
			db.Properties.Resource.PartitionKey.Version = utils.Int64(int64(partitionKeyVersion.(int)))
		}
	}

	if keys := expandAzureRmCosmosDbGremlinGraphUniqueKeys(d.Get("unique_key").(*pluginsdk.Set)); keys != nil {
		db.Properties.Resource.UniqueKeyPolicy = &cosmosdb.UniqueKeyPolicy{
			UniqueKeys: keys,
		}
	}

	if v, ok := d.GetOk("analytical_storage_ttl"); ok {
		db.Properties.Resource.AnalyticalStorageTtl = utils.Int64(int64(v.(int)))
	}

	if defaultTTL, hasDefaultTTL := d.GetOk("default_ttl"); hasDefaultTTL {
		db.Properties.Resource.DefaultTtl = utils.Int64(int64(defaultTTL.(int)))
	}

	err = client.GremlinResourcesCreateUpdateGremlinGraphThenPoll(ctx, *id, db)
	if err != nil {
		return fmt.Errorf("updating %q: %+v", id, err)
	}

	if common.HasThroughputChange(d) {
		throughputParameters := common.ExpandCosmosDBThroughputSettingsUpdateParameters(d)
		throughputFuture, err := client.GremlinResourcesUpdateGremlinGraphThroughput(ctx, *id, *throughputParameters)
		if err != nil {
			if response.WasNotFound(throughputFuture.HttpResponse) {
				return fmt.Errorf("setting Throughput for Cosmos Gremlin Graph %q (Account: %q, Database: %q): %+v - "+
					"If the graph has not been created with an initial throughput, you cannot configure it later.", id.GraphName, id.DatabaseAccountName, id.GremlinDatabaseName, err)
			}
		}

		if err := throughputFuture.Poller.PollUntilDone(); err != nil {
			return fmt.Errorf("waiting on ThroughputUpdate future for Cosmos Gremlin Graph %q (Account: %q, Database: %q): %+v", id.GraphName, id.DatabaseAccountName, id.GremlinDatabaseName, err)
		}
	}

	return resourceCosmosDbGremlinGraphRead(d, meta)
}

func resourceCosmosDbGremlinGraphRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient
	accountClient := meta.(*clients.Client).Cosmos.DatabaseClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseGraphID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GremlinResourcesGetGremlinGraph(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Error reading %q - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading %q: %+v", id, err)
	}

	d.Set("name", id.GraphName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("account_name", id.DatabaseAccountName)
	d.Set("database_name", id.GremlinDatabaseName)

	if model := resp.Model; model != nil {
		if graphProperties := model.Properties; graphProperties != nil {
			if props := graphProperties.Resource; props != nil {
				if pk := props.PartitionKey; pk != nil {
					if paths := pk.Paths; paths != nil {
						if len(*paths) > 1 {
							return fmt.Errorf("reading PartitionKey Paths, more than 1 returned")
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
						return fmt.Errorf("setting `index_policy`: %+v", err)
					}
				}

				if crp := props.ConflictResolutionPolicy; crp != nil {
					if err := d.Set("conflict_resolution_policy", common.FlattenCosmosDbConflictResolutionPolicy(crp)); err != nil {
						return fmt.Errorf("setting `conflict_resolution_policy`: %+v", err)
					}
				}

				if ukp := props.UniqueKeyPolicy; ukp != nil {
					if err := d.Set("unique_key", flattenCosmosGremlinGraphUniqueKeys(ukp.UniqueKeys)); err != nil {
						return fmt.Errorf("setting `unique_key`: %+v", err)
					}
				}

				if v := props.AnalyticalStorageTtl; v != nil {
					d.Set("analytical_storage_ttl", v)
				}

				if defaultTTL := props.DefaultTtl; defaultTTL != nil {
					d.Set("default_ttl", defaultTTL)
				}
			}
		}
	}
	accResp, err := accountClient.Get(ctx, id.ResourceGroupName, id.DatabaseAccountName)
	if err != nil {
		return fmt.Errorf("reading Cosmos Account %q : %+v", id.DatabaseAccountName, err)
	}
	if accResp.ID == nil || *accResp.ID == "" {
		return fmt.Errorf("cosmosDB Account %q (Resource Group %q) ID is empty or nil", id.DatabaseAccountName, id.ResourceGroupName)
	}

	if !isServerlessCapacityMode(accResp) {
		throughputResp, err := client.GremlinResourcesGetGremlinGraphThroughput(ctx, *id)
		if err != nil {
			if !response.WasNotFound(throughputResp.HttpResponse) {
				return fmt.Errorf("reading Throughput on Gremlin Graph %q (Account: %q, Database: %q) ID: %v", id.GraphName, id.DatabaseAccountName, id.GremlinDatabaseName, err)
			} else {
				d.Set("throughput", nil)
				d.Set("autoscale_settings", nil)
			}
		} else {
			common.SetResourceDataThroughputFromResponse(*throughputResp.Model, d)
		}
	}
	return nil
}

func resourceCosmosDbGremlinGraphDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseGraphID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.GremlinResourcesDeleteGremlinGraph(ctx, *id)
	if err != nil {
		if !response.WasNotFound(future.HttpResponse) {
			return fmt.Errorf("deleting Cosmos Gremlin Graph %q (Account: %q): %+v", id.GremlinDatabaseName, id.GraphName, err)
		}
	}

	if err := future.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("waiting on delete future for Comos Gremlin Graph %q (Account: %q): %+v", id.GremlinDatabaseName, id.DatabaseAccountName, err)
	}

	return nil
}

func expandAzureRmCosmosDbGrelinGraphIndexingPolicy(d *pluginsdk.ResourceData) *cosmosdb.IndexingPolicy {
	i := d.Get("index_policy").([]interface{})
	if len(i) == 0 || i[0] == nil {
		return nil
	}

	input := i[0].(map[string]interface{})
	indexingPolicy := cosmosdb.IndexingMode(strings.ToLower(input["indexing_mode"].(string)))
	policy := &cosmosdb.IndexingPolicy{
		IndexingMode:  &indexingPolicy,
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

func expandAzureRmCosmosDbGrelimGraphIncludedPath(input map[string]interface{}) *[]cosmosdb.IncludedPath {
	includedPath := input["included_paths"].(*pluginsdk.Set).List()
	paths := make([]cosmosdb.IncludedPath, len(includedPath))

	for i, pathConfig := range includedPath {
		attrs := pathConfig.(string)
		path := cosmosdb.IncludedPath{
			Path: utils.String(attrs),
		}
		paths[i] = path
	}

	return &paths
}

func expandAzureRmCosmosDbGremlinGraphExcludedPath(input map[string]interface{}) *[]cosmosdb.ExcludedPath {
	excludedPath := input["excluded_paths"].(*pluginsdk.Set).List()
	paths := make([]cosmosdb.ExcludedPath, len(excludedPath))

	for i, pathConfig := range excludedPath {
		attrs := pathConfig.(string)
		path := cosmosdb.ExcludedPath{
			Path: utils.String(attrs),
		}
		paths[i] = path
	}

	return &paths
}

func expandAzureRmCosmosDbGremlinGraphUniqueKeys(s *pluginsdk.Set) *[]cosmosdb.UniqueKey {
	i := s.List()
	if len(i) == 0 || i[0] == nil {
		return nil
	}

	keys := make([]cosmosdb.UniqueKey, 0)
	for _, k := range i {
		key := k.(map[string]interface{})

		paths := key["paths"].(*pluginsdk.Set).List()
		if len(paths) == 0 {
			continue
		}

		keys = append(keys, cosmosdb.UniqueKey{
			Paths: utils.ExpandStringSlice(paths),
		})
	}

	return &keys
}

func flattenAzureRmCosmosDBGremlinGraphIndexingPolicy(input *cosmosdb.IndexingPolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	indexPolicy := make(map[string]interface{})

	indexPolicy["automatic"] = input.Automatic
	indexPolicy["indexing_mode"] = input.IndexingMode
	indexPolicy["included_paths"] = pluginsdk.NewSet(pluginsdk.HashString, flattenAzureRmCosmosDBGremlinGraphIncludedPaths(input.IncludedPaths))
	indexPolicy["excluded_paths"] = pluginsdk.NewSet(pluginsdk.HashString, flattenAzureRmCosmosDBGremlinGraphExcludedPaths(input.ExcludedPaths))
	indexPolicy["composite_index"] = common.FlattenCosmosDBIndexingPolicyCompositeIndexes(input.CompositeIndexes)
	indexPolicy["spatial_index"] = common.FlattenCosmosDBIndexingPolicySpatialIndexes(input.SpatialIndexes)

	return []interface{}{indexPolicy}
}

func flattenAzureRmCosmosDBGremlinGraphIncludedPaths(input *[]cosmosdb.IncludedPath) []interface{} {
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

func flattenAzureRmCosmosDBGremlinGraphExcludedPaths(input *[]cosmosdb.ExcludedPath) []interface{} {
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

func flattenCosmosGremlinGraphUniqueKeys(keys *[]cosmosdb.UniqueKey) *[]map[string]interface{} {
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
