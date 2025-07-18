// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedredis

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagedRedisDatabaseResource struct{}

var _ sdk.ResourceWithUpdate = ManagedRedisDatabaseResource{}

type ManagedRedisDatabaseResourceModel struct {
	Name                            string        `tfschema:"name"`
	ClusterId                       string        `tfschema:"cluster_id"`
	AccessKeysAuthenticationEnabled bool          `tfschema:"access_keys_authentication_enabled"`
	ClientProtocol                  string        `tfschema:"client_protocol"`
	ClusteringPolicy                string        `tfschema:"clustering_policy"`
	EvictionPolicy                  string        `tfschema:"eviction_policy"`
	LinkedDatabaseGroupNickname     string        `tfschema:"linked_database_group_nickname"`
	LinkedDatabaseId                []string      `tfschema:"linked_database_id"`
	Module                          []ModuleModel `tfschema:"module"`
	Port                            int64         `tfschema:"port"`
	PrimaryAccessKey                string        `tfschema:"primary_access_key"`
	SecondaryAccessKey              string        `tfschema:"secondary_access_key"`
}

type ModuleModel struct {
	Name    string `tfschema:"name"`
	Args    string `tfschema:"args"`
	Version string `tfschema:"version"`
}

func (r ManagedRedisDatabaseResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "default",
			ValidateFunc: validate.ManagedRedisDatabaseName,
		},

		"cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: redisenterprise.ValidateRedisEnterpriseID,
		},

		"access_keys_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"client_protocol": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(redisenterprise.ProtocolEncrypted),
			ValidateFunc: validation.StringInSlice(redisenterprise.PossibleValuesForProtocol(), false),
		},

		"clustering_policy": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(redisenterprise.ClusteringPolicyOSSCluster),
			ValidateFunc: validation.StringInSlice(redisenterprise.PossibleValuesForClusteringPolicy(), false),
		},

		"eviction_policy": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(redisenterprise.EvictionPolicyVolatileLRU),
			ValidateFunc: validation.StringInSlice(redisenterprise.PossibleValuesForEvictionPolicy(), false),
		},

		"linked_database_group_nickname": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"linked_database_id"},
		},

		"linked_database_id": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			MaxItems: 5,
			Set:      pluginsdk.HashString,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: databases.ValidateDatabaseID,
			},
		},

		"module": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 4,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							"RedisBloom",
							"RedisTimeSeries",
							"RediSearch",
							"RedisJSON",
						}, false),
					},

					"args": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  "",
					},

					"version": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"port": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      10000,
			ValidateFunc: validation.IntBetween(0, 65353),
		},
	}
}

func (r ManagedRedisDatabaseResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"primary_access_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_access_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},
	}
}

func (r ManagedRedisDatabaseResource) ModelObject() interface{} {
	return &ManagedRedisDatabaseResourceModel{}
}

func (r ManagedRedisDatabaseResource) ResourceType() string {
	return "azurerm_managed_redis_database"
}

func (r ManagedRedisDatabaseResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return databases.ValidateDatabaseID
}

func (r ManagedRedisDatabaseResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.ManagedRedis.DatabaseClient

			var model ManagedRedisDatabaseResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			clusterId, err := redisenterprise.ParseRedisEnterpriseID(model.ClusterId)
			if err != nil {
				return fmt.Errorf("parsing `cluster_id`: %+v", err)
			}

			id := databases.NewDatabaseID(subscriptionId, clusterId.ResourceGroupName, clusterId.RedisEnterpriseName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			clusteringPolicy := databases.ClusteringPolicy(model.ClusteringPolicy)
			evictionPolicy := databases.EvictionPolicy(model.EvictionPolicy)
			protocol := databases.Protocol(model.ClientProtocol)

			accessKeysAuth := databases.AccessKeysAuthenticationDisabled
			if model.AccessKeysAuthenticationEnabled {
				accessKeysAuth = databases.AccessKeysAuthenticationEnabled
			}

			linkedDatabase, err := expandArmGeoLinkedDatabase(model.LinkedDatabaseId, id.ID(), model.LinkedDatabaseGroupNickname)
			if err != nil {
				return fmt.Errorf("Setting geo database for database %s error: %+v", id.ID(), err)
			}

			isGeoEnabled := false
			if linkedDatabase != nil {
				isGeoEnabled = true
			}
			module, err := expandArmDatabaseModuleArray(model.Module, isGeoEnabled)
			if err != nil {
				return fmt.Errorf("setting module error: %+v", err)
			}

			parameters := databases.Database{
				Properties: &databases.DatabaseProperties{
					AccessKeysAuthentication: &accessKeysAuth,
					ClientProtocol:           &protocol,
					ClusteringPolicy:         &clusteringPolicy,
					EvictionPolicy:           &evictionPolicy,
					Port:                     utils.Int64(model.Port),
					GeoReplication:           linkedDatabase,
					Modules:                  module,
				},
			}

			future, err := client.Create(ctx, id, parameters)
			if err != nil {
				// Need to check if this was due to the cluster having the wrong sku
				if strings.Contains(err.Error(), "The value of the parameter 'properties.modules' is invalid") {
					clusterClient := metadata.Client.ManagedRedis.Client
					resp, err := clusterClient.Get(ctx, *clusterId)
					if err != nil {
						return fmt.Errorf("retrieving %s: %+v", *clusterId, err)
					}

					if resp.Model != nil && strings.Contains(strings.ToLower(string(resp.Model.Sku.Name)), "flash") {
						return fmt.Errorf("creating a Managed Redis Database with modules in a Managed Redis Cluster that has an incompatible Flash SKU type %q - please remove the Managed Redis Database modules or change the Managed Redis Cluster SKU type %s", string(resp.Model.Sku.Name), id)
					}
				}

				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err := future.Poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagedRedisDatabaseResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedRedis.DatabaseClient

			id, err := databases.ParseDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ManagedRedisDatabaseResourceModel{
				Name: id.DatabaseName,
			}

			clusterId := redisenterprise.NewRedisEnterpriseID(id.SubscriptionId, id.ResourceGroupName, id.RedisEnterpriseName)
			state.ClusterId = clusterId.ID()

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.AccessKeysAuthenticationEnabled = strings.EqualFold(string(pointer.From(props.AccessKeysAuthentication)), string(databases.AccessKeysAuthenticationEnabled))
					state.ClientProtocol = string(pointer.From(props.ClientProtocol))
					state.ClusteringPolicy = string(pointer.From(props.ClusteringPolicy))
					state.EvictionPolicy = string(pointer.From(props.EvictionPolicy))
					state.Port = pointer.From(props.Port)

					if geoProps := props.GeoReplication; geoProps != nil {
						if geoProps.GroupNickname != nil {
							state.LinkedDatabaseGroupNickname = pointer.From(geoProps.GroupNickname)
						}
						state.LinkedDatabaseId = flattenArmGeoLinkedDatabase(geoProps.LinkedDatabases)
					}

					state.Module = flattenArmDatabaseModuleArray(props.Modules)

					if state.AccessKeysAuthenticationEnabled {
						keysResp, err := client.ListKeys(ctx, *id)
						if err != nil {
							return fmt.Errorf("listing keys for %s: %+v", *id, err)
						}
						if keysModel := keysResp.Model; keysModel != nil {
							state.PrimaryAccessKey = pointer.From(keysResp.Model.PrimaryKey)
							state.SecondaryAccessKey = pointer.From(keysResp.Model.SecondaryKey)
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagedRedisDatabaseResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedRedis.DatabaseClient

			id, err := databases.ParseDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagedRedisDatabaseResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil || existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model or properties was nil", *id)
			}

			parameters := databases.Database{
				Properties: existing.Model.Properties,
			}

			linkedDatabase, err := expandArmGeoLinkedDatabase(model.LinkedDatabaseId, id.ID(), model.LinkedDatabaseGroupNickname)
			if err != nil {
				return fmt.Errorf("Setting geo database for database %s error: %+v", id.ID(), err)
			}

			isGeoEnabled := false
			if linkedDatabase != nil {
				isGeoEnabled = true
			}

			if metadata.ResourceData.HasChange("access_keys_authentication_enabled") {
				accessKeysAuth := databases.AccessKeysAuthenticationDisabled
				if model.AccessKeysAuthenticationEnabled {
					accessKeysAuth = databases.AccessKeysAuthenticationEnabled
				}
				parameters.Properties.AccessKeysAuthentication = &accessKeysAuth
			}

			if metadata.ResourceData.HasChange("module") {
				module, err := expandArmDatabaseModuleArray(model.Module, isGeoEnabled)
				if err != nil {
					return fmt.Errorf("setting module error: %+v", err)
				}
				parameters.Properties.Modules = module
			}

			if metadata.ResourceData.HasChange("linked_database_id") {
				oldItems, newItems := metadata.ResourceData.GetChange("linked_database_id")
				isForceUnlink, data := forceUnlinkItems(oldItems.(*pluginsdk.Set).List(), newItems.(*pluginsdk.Set).List())
				if isForceUnlink {
					if err := forceUnlinkDatabase(&metadata, data); err != nil {
						return fmt.Errorf("unlinking database error: %+v", err)
					}
				}
				parameters.Properties.GeoReplication = linkedDatabase
			}

			if err := client.CreateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ManagedRedisDatabaseResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedRedis.DatabaseClient
			clusterClient := metadata.Client.ManagedRedis.Client

			id, err := databases.ParseDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			dbId := databases.NewDatabaseID(id.SubscriptionId, id.ResourceGroupName, id.RedisEnterpriseName, id.DatabaseName)
			clusterId := redisenterprise.NewRedisEnterpriseID(id.SubscriptionId, id.ResourceGroupName, id.RedisEnterpriseName)

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			// can't use DeleteThenPoll since cluster deletion also deletes the default database, which will cause db deletion failure
			stateConf := &pluginsdk.StateChangeConf{
				Pending:                   []string{"found"},
				Target:                    []string{"clusterNotFound", "dbNotFound"},
				Refresh:                   redisEnterpriseDatabaseDeleteRefreshFunc(ctx, client, clusterClient, clusterId, dbId),
				ContinuousTargetOccurence: 3,
				Timeout:                   metadata.ResourceData.Timeout(pluginsdk.TimeoutDelete),
			}

			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for deletion %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandArmDatabaseModuleArray(input []ModuleModel, isGeoEnabled bool) (*[]databases.Module, error) {
	results := make([]databases.Module, 0)

	for _, item := range input {
		if item.Name != "RediSearch" && item.Name != "RedisJSON" && isGeoEnabled {
			return nil, fmt.Errorf("Only RediSearch and RedisJSON modules are allowed with geo-replication")
		}
		results = append(results, databases.Module{
			Name: item.Name,
			Args: utils.String(item.Args),
		})
	}
	return &results, nil
}

func flattenArmDatabaseModuleArray(input *[]databases.Module) []ModuleModel {
	results := make([]ModuleModel, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		args := ""
		if item.Args != nil {
			args = *item.Args
			if strings.EqualFold(args, "PARTITIONS AUTO") {
				args = ""
			}
		}

		var version string
		if item.Version != nil {
			version = *item.Version
		}

		results = append(results, ModuleModel{
			Name:    item.Name,
			Args:    args,
			Version: version,
		})
	}

	return results
}

func forceUnlinkDatabase(meta *sdk.ResourceMetaData, unlinkedDbRaw []string) error {
	client := meta.Client.ManagedRedis.DatabaseClient
	ctx, cancel := timeouts.ForUpdate(meta.Client.StopContext, meta.ResourceData)
	defer cancel()
	log.Printf("[INFO] Preparing to unlink a linked database")

	id, err := databases.ParseDatabaseID(meta.ResourceData.Id())
	if err != nil {
		return err
	}

	parameters := databases.ForceUnlinkParameters{
		Ids: unlinkedDbRaw,
	}

	if err := client.ForceUnlinkThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("force unlinking from database %s error: %+v", id, err)
	}

	return nil
}

func expandArmGeoLinkedDatabase(inputId []string, parentDBId string, inputGeoName string) (*databases.DatabasePropertiesGeoReplication, error) {
	idList := make([]databases.LinkedDatabase, 0)
	if len(inputId) == 0 {
		return nil, nil
	}
	isParentDbIncluded := false

	for _, id := range inputId {
		if id == parentDBId {
			isParentDbIncluded = true
		}
		idList = append(idList, databases.LinkedDatabase{
			Id: utils.String(id),
		})
	}
	if isParentDbIncluded {
		return &databases.DatabasePropertiesGeoReplication{
			LinkedDatabases: &idList,
			GroupNickname:   utils.String(inputGeoName),
		}, nil
	}

	return nil, fmt.Errorf("linked database list must include database ID: %s", parentDBId)
}

func flattenArmGeoLinkedDatabase(inputDB *[]databases.LinkedDatabase) []string {
	results := make([]string, 0)

	if inputDB == nil {
		return results
	}

	for _, item := range *inputDB {
		if item.Id != nil {
			results = append(results, *item.Id)
		}
	}
	return results
}

func forceUnlinkItems(oldItemList []interface{}, newItemList []interface{}) (bool, []string) {
	newItems := make(map[string]bool)
	forceUnlinkList := make([]string, 0)
	for _, newItem := range newItemList {
		newItems[newItem.(string)] = true
	}

	for _, oldItem := range oldItemList {
		if !newItems[oldItem.(string)] {
			forceUnlinkList = append(forceUnlinkList, oldItem.(string))
		}
	}
	if len(forceUnlinkList) > 0 {
		return true, forceUnlinkList
	}
	return false, nil
}

func redisEnterpriseDatabaseDeleteRefreshFunc(ctx context.Context, databaseClient *databases.DatabasesClient, clusterClient *redisenterprise.RedisEnterpriseClient, clusterId redisenterprise.RedisEnterpriseId, databaseId databases.DatabaseId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		cluster, err := clusterClient.Get(ctx, clusterId)
		if err != nil {
			if response.WasNotFound(cluster.HttpResponse) {
				return "clusterNotFound", "clusterNotFound", nil
			}
		}
		db, err := databaseClient.Get(ctx, databaseId)
		if err != nil {
			if response.WasNotFound(db.HttpResponse) {
				return "dbNotFound", "dbNotFound", nil
			}
		}
		return db, "found", nil
	}
}
