// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"context"
	"fmt"
	"log"
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

func resourceCosmosDbSQLContainer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCosmosDbSQLContainerCreate,
		Read:   resourceCosmosDbSQLContainerRead,
		Update: resourceCosmosDbSQLContainerUpdate,
		Delete: resourceCosmosDbSQLContainerDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SqlContainerID(id)
			return err
		}),

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

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			// The analytical_storage_ttl cannot be changed back once enabled on an existing container. -> we need ForceNew
			pluginsdk.ForceNewIfChange("analytical_storage_ttl", func(ctx context.Context, old, new, _ interface{}) bool {
				return (old.(int) == -1 || old.(int) > 0) && new.(int) == 0
			}),
		),
	}
}

func resourceCosmosDbSQLContainerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := cosmosdb.NewContainerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("database_name").(string), d.Get("name").(string))
	partitionkeypaths := d.Get("partition_key_path").(string)

	existing, err := client.SqlResourcesGetSqlContainer(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_cosmosdb_sql_container", id.ID())
	}

	indexingPolicy := common.ExpandAzureRmCosmosDbIndexingPolicy(d)
	err = common.ValidateAzureRmCosmosDbIndexingPolicy(indexingPolicy)
	if err != nil {
		return fmt.Errorf("generating indexing policy for %s", id)
	}

	db := cosmosdb.SqlContainerCreateUpdateParameters{
		Properties: cosmosdb.SqlContainerCreateUpdateProperties{
			Resource: cosmosdb.SqlContainerResource{
				Id:                       id.ContainerName,
				IndexingPolicy:           indexingPolicy,
				ConflictResolutionPolicy: common.ExpandCosmosDbConflicResolutionPolicy(d.Get("conflict_resolution_policy").([]interface{})),
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

	if keys := expandCosmosSQLContainerUniqueKeys(d.Get("unique_key").(*pluginsdk.Set)); keys != nil {
		db.Properties.Resource.UniqueKeyPolicy = &cosmosdb.UniqueKeyPolicy{
			UniqueKeys: keys,
		}
	}

	if analyticalStorageTTL, ok := d.GetOk("analytical_storage_ttl"); ok {
		db.Properties.Resource.AnalyticalStorageTtl = utils.Int64(int64(analyticalStorageTTL.(int)))
	}

	if defaultTTL, hasTTL := d.GetOk("default_ttl"); hasTTL {
		db.Properties.Resource.DefaultTtl = utils.Int64(int64(defaultTTL.(int)))
	}

	if throughput, hasThroughput := d.GetOk("throughput"); hasThroughput {
		if throughput != 0 {
			db.Properties.Options.Throughput = common.ConvertThroughputFromResourceData(throughput)
		}
	}

	if _, hasAutoscaleSettings := d.GetOk("autoscale_settings"); hasAutoscaleSettings {
		db.Properties.Options.AutoScaleSettings = common.ExpandCosmosDbAutoscaleSettings(d)
	}

	err = client.SqlResourcesCreateUpdateSqlContainerThenPoll(ctx, id, db)
	if err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosDbSQLContainerRead(d, meta)
}

func resourceCosmosDbSQLContainerUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseContainerID(d.Id())
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

	db := cosmosdb.SqlContainerCreateUpdateParameters{
		Properties: cosmosdb.SqlContainerCreateUpdateProperties{
			Resource: cosmosdb.SqlContainerResource{
				Id:             id.ContainerName,
				IndexingPolicy: indexingPolicy,
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

	if keys := expandCosmosSQLContainerUniqueKeys(d.Get("unique_key").(*pluginsdk.Set)); keys != nil {
		db.Properties.Resource.UniqueKeyPolicy = &cosmosdb.UniqueKeyPolicy{
			UniqueKeys: keys,
		}
	}

	if analyticalStorageTTL, ok := d.GetOk("analytical_storage_ttl"); ok {
		db.Properties.Resource.AnalyticalStorageTtl = utils.Int64(int64(analyticalStorageTTL.(int)))
	}

	if defaultTTL, hasTTL := d.GetOk("default_ttl"); hasTTL {
		db.Properties.Resource.DefaultTtl = utils.Int64(int64(defaultTTL.(int)))
	}

	err = client.SqlResourcesCreateUpdateSqlContainerThenPoll(ctx, *id, db)
	if err != nil {
		return fmt.Errorf("updating %q: %+v", id, err)
	}

	if common.HasThroughputChange(d) {
		throughputParameters := common.ExpandCosmosDBThroughputSettingsUpdateParameters(d)
		throughputFuture, err := client.SqlResourcesUpdateSqlContainerThroughput(ctx, *id, *throughputParameters)
		if err != nil {
			if response.WasNotFound(throughputFuture.HttpResponse) {
				return fmt.Errorf("setting Throughput for Cosmos SQL Container %q (Account: %q, Database: %q): %+v - "+
					"If the collection has not been created with an initial throughput, you cannot configure it later", id.ContainerName, id.DatabaseAccountName, id.SqlDatabaseName, err)
			}
		}

		if err := throughputFuture.Poller.PollUntilDone(); err != nil {
			return fmt.Errorf("waiting on ThroughputUpdate future for Cosmos Container %q (Account: %q, Database: %q): %+v", id.ContainerName, id.DatabaseAccountName, id.SqlDatabaseName, err)
		}
	}

	return resourceCosmosDbSQLContainerRead(d, meta)
}

func resourceCosmosDbSQLContainerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient
	accountClient := meta.(*clients.Client).Cosmos.DatabaseClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseContainerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.SqlResourcesGetSqlContainer(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Error reading %q - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading %q: %+v", id, err)
	}

	d.Set("name", id.ContainerName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("account_name", id.DatabaseAccountName)
	d.Set("database_name", id.SqlDatabaseName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
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

				if analyticalStorageTTL := res.AnalyticalStorageTtl; analyticalStorageTTL != nil {
					d.Set("analytical_storage_ttl", analyticalStorageTTL)
				}

				if defaultTTL := res.DefaultTtl; defaultTTL != nil {
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
	}

	accResp, err := accountClient.Get(ctx, id.ResourceGroupName, id.DatabaseAccountName)
	if err != nil {
		return fmt.Errorf("reading CosmosDB Account %q (Resource Group %q): %+v", id.DatabaseAccountName, id.ResourceGroupName, err)
	}

	if accResp.ID == nil || *accResp.ID == "" {
		return fmt.Errorf("cosmosDB Account %q (Resource Group %q) ID is empty or nil", id.DatabaseAccountName, id.ResourceGroupName)
	}

	// if the cosmos account is serverless calling the get throughput api would yield an error
	if !isServerlessCapacityMode(accResp) {
		throughputResp, err := client.SqlResourcesGetSqlContainerThroughput(ctx, *id)
		if err != nil {
			if !response.WasNotFound(throughputResp.HttpResponse) {
				return fmt.Errorf("reading Throughput on Cosmos SQL Container %s (Account: %q, Database: %q) ID: %v", id.ContainerName, id.DatabaseAccountName, id.SqlDatabaseName, err)
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

func resourceCosmosDbSQLContainerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseContainerID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.SqlResourcesDeleteSqlContainer(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting Cosmos SQL Container %q (Account: %q): %+v", id.SqlDatabaseName, id.ContainerName, err)
	}

	if err := future.Poller.PollUntilDone(); err != nil {
		if !response.WasNotFound(future.HttpResponse) {
			return fmt.Errorf("deleting Cosmos SQL Container %q (Account: %q): %+v", id.SqlDatabaseName, id.ContainerName, err)
		}
	}

	return nil
}

func expandCosmosSQLContainerUniqueKeys(s *pluginsdk.Set) *[]cosmosdb.UniqueKey {
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

func flattenCosmosSQLContainerUniqueKeys(keys *[]cosmosdb.UniqueKey) *[]map[string]interface{} {
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
