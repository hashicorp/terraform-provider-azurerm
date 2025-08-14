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
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/databaselink"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagedRedisDatabaseResource struct{}

var (
	_ sdk.ResourceWithCustomizeDiff = ManagedRedisDatabaseResource{}
	_ sdk.ResourceWithUpdate        = ManagedRedisDatabaseResource{}
)

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
			ValidateFunc: validate.ManagedRedisDatabaseGeoreplicationGroupName,
		},

		"linked_database_id": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			MaxItems: 5,
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
			client := metadata.Client.ManagedRedis.DatabaseClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ManagedRedisDatabaseResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			clusterId, err := redisenterprise.ParseRedisEnterpriseID(model.ClusterId)
			if err != nil {
				return err
			}

			id := databases.NewDatabaseID(subscriptionId, clusterId.ResourceGroupName, clusterId.RedisEnterpriseName, model.Name)

			if err := validate.ValidateLinkedDatabaseIncludesSelf(model.LinkedDatabaseId, id.ID()); err != nil {
				return err
			}

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			accessKeysAuth := databases.AccessKeysAuthenticationDisabled
			if model.AccessKeysAuthenticationEnabled {
				accessKeysAuth = databases.AccessKeysAuthenticationEnabled
			}

			linkedDatabase, err := expandArmGeoLinkedDatabase(model.LinkedDatabaseId, model.LinkedDatabaseGroupNickname)
			if err != nil {
				return fmt.Errorf("expanding `linked_database_group_nickname` or `linked_database_id`: %+v", err)
			}

			isGeoEnabled := false
			if linkedDatabase != nil {
				isGeoEnabled = true
			}
			module, err := expandArmDatabaseModuleArray(model.Module, isGeoEnabled)
			if err != nil {
				return fmt.Errorf("expanding `module`: %+v", err)
			}

			parameters := databases.Database{
				Properties: &databases.DatabaseProperties{
					AccessKeysAuthentication: &accessKeysAuth,
					ClientProtocol:           pointer.To(databases.Protocol(model.ClientProtocol)),
					ClusteringPolicy:         pointer.To(databases.ClusteringPolicy(model.ClusteringPolicy)),
					EvictionPolicy:           pointer.To(databases.EvictionPolicy(model.EvictionPolicy)),
					Port:                     pointer.To(model.Port),
					GeoReplication:           linkedDatabase,
					Modules:                  module,
				},
			}

			future, err := client.Create(ctx, id, parameters)
			if err != nil {
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
						state.LinkedDatabaseGroupNickname = pointer.From(geoProps.GroupNickname)
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
				return fmt.Errorf("decoding: %+v", err)
			}

			if err := validate.ValidateLinkedDatabaseIncludesSelf(model.LinkedDatabaseId, id.ID()); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			param := existing.Model
			databaseUpdateNeeded := false

			if metadata.ResourceData.HasChange("linked_database_id") {
				oldItems, newItems := metadata.ResourceData.GetChange("linked_database_id")
				if err := forceUnlinkDatabase(ctx, &metadata, oldItems, newItems); err != nil {
					return fmt.Errorf("force unlinking database error: %+v", err)
				}
				if err := forceLinkDatabase(ctx, &metadata, oldItems, newItems); err != nil {
					return fmt.Errorf("force linking database error: %+v", err)
				}
			}

			if metadata.ResourceData.HasChange("access_keys_authentication_enabled") {
				param.Properties.AccessKeysAuthentication = pointer.To(databases.AccessKeysAuthenticationDisabled)
				if model.AccessKeysAuthenticationEnabled {
					param.Properties.AccessKeysAuthentication = pointer.To(databases.AccessKeysAuthenticationEnabled)
				}
				databaseUpdateNeeded = true
			}

			if databaseUpdateNeeded {
				// Oddly this SDK does not have a CreateOrUpdate method. Despite the name, Create uses PUT. Update / PATCH
				// method cannot be used because it has a bug where accessKeysAuthentication update is not yet implemented.
				if err := client.CreateThenPoll(ctx, *id, *param); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
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
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			// can't use DeleteThenPoll since cluster deletion also deletes the default database, which will cause db deletion failure
			deletePoller := custompollers.NewDatabaseDeletePoller(client, clusterClient, dbId, clusterId)
			poller := pollers.NewPoller(deletePoller, 15*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)

			timeoutContext, cancel := context.WithTimeout(ctx, metadata.ResourceData.Timeout(pluginsdk.TimeoutDelete))
			defer cancel()

			if err := poller.PollUntilDone(timeoutContext); err != nil {
				return fmt.Errorf("waiting for deletion %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r ManagedRedisDatabaseResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}

			var model ManagedRedisDatabaseResourceModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return err
			}

			isGeoEnabled := len(model.LinkedDatabaseId) > 0

			if isGeoEnabled {
				var currentDatabaseId string
				switch {
				case metadata.ResourceData != nil && metadata.ResourceData.Id() != "":
					currentDatabaseId = metadata.ResourceData.Id()
				case model.ClusterId != "" && model.Name != "":
					clusterId, err := redisenterprise.ParseRedisEnterpriseID(model.ClusterId)
					if err != nil {
						return fmt.Errorf("parsing cluster ID: %+v", err)
					}
					currentDatabaseId = databases.NewDatabaseID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.RedisEnterpriseName, model.Name).ID()
				default:
					return nil
				}

				isCurrentDbIncluded := false
				for _, id := range model.LinkedDatabaseId {
					if id == currentDatabaseId {
						isCurrentDbIncluded = true
						break
					}
				}

				if !isCurrentDbIncluded {
					return fmt.Errorf("linked database list must include the current database ID: %s", currentDatabaseId)
				}

				for _, module := range model.Module {
					if module.Name != "RediSearch" && module.Name != "RedisJSON" {
						return fmt.Errorf("Only `RediSearch` and `RedisJSON` modules are allowed with geo-replication")
					}
				}
			}

			return nil
		},
	}
}

func expandArmDatabaseModuleArray(input []ModuleModel, isGeoEnabled bool) (*[]databases.Module, error) {
	results := make([]databases.Module, 0)

	for _, item := range input {
		if item.Name != "RediSearch" && item.Name != "RedisJSON" && isGeoEnabled {
			return nil, fmt.Errorf("Only `RediSearch` and `RedisJSON` modules are allowed with geo-replication")
		}
		results = append(results, databases.Module{
			Name: item.Name,
			Args: pointer.To(item.Args),
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
		results = append(results, ModuleModel{
			Name:    item.Name,
			Args:    pointer.From(item.Args),
			Version: pointer.From(item.Version),
		})
	}

	return results
}

func expandArmGeoLinkedDatabase(inputId []string, inputGeoName string) (*databases.DatabasePropertiesGeoReplication, error) {
	idList := make([]databases.LinkedDatabase, 0)
	if len(inputId) == 0 {
		return nil, nil
	}

	for _, id := range inputId {
		idList = append(idList, databases.LinkedDatabase{
			Id: pointer.To(id),
		})
	}

	return &databases.DatabasePropertiesGeoReplication{
		LinkedDatabases: &idList,
		GroupNickname:   pointer.To(inputGeoName),
	}, nil
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

func forceUnlinkDatabase(ctx context.Context, meta *sdk.ResourceMetaData, oldItems, newItems interface{}) error {
	isForceUnlinkNeeded, data := databaselink.ForceUnlinkItems(oldItems.(*pluginsdk.Set).List(), newItems.(*pluginsdk.Set).List())
	if isForceUnlinkNeeded {
		client := meta.Client.ManagedRedis.DatabaseClient
		log.Printf("[INFO] Preparing to unlink a linked database")

		id, err := databases.ParseDatabaseID(meta.ResourceData.Id())
		if err != nil {
			return err
		}

		parameters := databases.ForceUnlinkParameters{
			Ids: data,
		}

		if err := client.ForceUnlinkThenPoll(ctx, *id, parameters); err != nil {
			return fmt.Errorf("force unlinking from database %s error: %+v", id, err)
		}
	}
	return nil
}

func forceLinkDatabase(ctx context.Context, meta *sdk.ResourceMetaData, oldItems, newItems interface{}) error {
	if databaselink.ForceLinkNeeded(oldItems.(*pluginsdk.Set).List(), newItems.(*pluginsdk.Set).List()) {
		client := meta.Client.ManagedRedis.DatabaseClient
		log.Printf("[INFO] Preparing to link to a replication group")

		id, err := databases.ParseDatabaseID(meta.ResourceData.Id())
		if err != nil {
			return err
		}

		var model ManagedRedisDatabaseResourceModel
		if err := meta.Decode(&model); err != nil {
			return fmt.Errorf("decoding model: %+v", err)
		}

		linkedDatabases := make([]databases.LinkedDatabase, 0)
		for _, item := range model.LinkedDatabaseId {
			linkedDatabases = append(linkedDatabases, databases.LinkedDatabase{
				Id: pointer.To(item),
			})
		}

		parameters := databases.ForceLinkParameters{
			GeoReplication: databases.ForceLinkParametersGeoReplication{
				GroupNickname:   pointer.To(model.LinkedDatabaseGroupNickname),
				LinkedDatabases: &linkedDatabases,
			},
		}

		if err := client.ForceLinkToReplicationGroupThenPoll(ctx, *id, parameters); err != nil {
			return fmt.Errorf("force linking to replication group %s error: %+v", id, err)
		}
	}
	return nil
}
