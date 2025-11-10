// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedredis

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// Azure Managed Redis (AMR) consists of two ARM resource types: cluster and database. Cluster is where compute, load
// balancer, network and other infrastructure is setup. Database refers to the Redis instance / process itself. Database
// is a child of cluster with 1-1 mapping.
//
// Database was its own resource in the deprecated redis_enterprise resource, but intentionally included here to improve
// UX. There were cases where users not aware of database and expect Redis Enterprise cluster to work by itself.
//
// There might be a plan to support multiple databases in the future, in which case we will implement an
// azurerm_managed_redis_custom_db resource.

type ManagedRedisResource struct{}

var _ sdk.ResourceWithUpdate = ManagedRedisResource{}

type ManagedRedisResourceModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`

	Location string `tfschema:"location"`

	SkuName string `tfschema:"sku_name"`

	CustomerManagedKey      []CustomerManagedKeyModel                  `tfschema:"customer_managed_key"`
	DefaultDatabase         []DefaultDatabaseModel                     `tfschema:"default_database"`
	HighAvailabilityEnabled bool                                       `tfschema:"high_availability_enabled"`
	Identity                []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	PublicNetworkAccess     string                                     `tfschema:"public_network_access"`
	Tags                    map[string]string                          `tfschema:"tags"`

	Hostname string `tfschema:"hostname"`
}

type CustomerManagedKeyModel struct {
	KeyVaultKeyId          string `tfschema:"key_vault_key_id"`
	UserAssignedIdentityId string `tfschema:"user_assigned_identity_id"`
}

type DefaultDatabaseModel struct {
	AccessKeysAuthenticationEnabled bool          `tfschema:"access_keys_authentication_enabled"`
	ClientProtocol                  string        `tfschema:"client_protocol"`
	ClusteringPolicy                string        `tfschema:"clustering_policy"`
	EvictionPolicy                  string        `tfschema:"eviction_policy"`
	GeoReplicationGroupName         string        `tfschema:"geo_replication_group_name"`
	Module                          []ModuleModel `tfschema:"module"`

	Port               int64  `tfschema:"port"`
	PrimaryAccessKey   string `tfschema:"primary_access_key"`
	SecondaryAccessKey string `tfschema:"secondary_access_key"`
}

type ModuleModel struct {
	Name    string `tfschema:"name"`
	Args    string `tfschema:"args"`
	Version string `tfschema:"version"`
}

const defaultDatabaseName = "default"

func (r ManagedRedisResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ManagedRedisClusterName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(validate.PossibleValuesForSkuName(), false),
		},

		"customer_managed_key": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
					},

					"user_assigned_identity_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
				},
			},
		},

		"default_database": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"access_keys_authentication_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"client_protocol": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(redisenterprise.ProtocolEncrypted),
						ValidateFunc: validation.StringInSlice(redisenterprise.PossibleValuesForProtocol(), false),
					},

					"clustering_policy": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(redisenterprise.ClusteringPolicyOSSCluster),
						ValidateFunc: validation.StringInSlice(redisenterprise.PossibleValuesForClusteringPolicy(), false),
					},

					"eviction_policy": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(redisenterprise.EvictionPolicyVolatileLRU),
						ValidateFunc: validation.StringInSlice(redisenterprise.PossibleValuesForEvictionPolicy(), false),
					},

					"geo_replication_group_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.ManagedRedisDatabaseGeoreplicationGroupName,
					},

					"module": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 4,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
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
								},

								"version": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"port": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

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
				},
			},
		},

		"high_availability_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  true,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"public_network_access": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      redisenterprise.PublicNetworkAccessEnabled,
			ValidateFunc: validation.StringInSlice(redisenterprise.PossibleValuesForPublicNetworkAccess(), false),
		},

		"tags": commonschema.Tags(),
	}
}

func (r ManagedRedisResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"hostname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ManagedRedisResource) ModelObject() interface{} {
	return &ManagedRedisResourceModel{}
}

func (r ManagedRedisResource) ResourceType() string {
	return "azurerm_managed_redis"
}

func (r ManagedRedisResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return redisenterprise.ValidateRedisEnterpriseID
}

func (r ManagedRedisResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 45 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			clusterClient := metadata.Client.ManagedRedis.Client
			dbClient := metadata.Client.ManagedRedis.DatabaseClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ManagedRedisResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			clusterId := redisenterprise.NewRedisEnterpriseID(subscriptionId, model.ResourceGroupName, model.Name)

			existingCluster, err := clusterClient.Get(ctx, clusterId)
			if err != nil {
				if !response.WasNotFound(existingCluster.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", clusterId, err)
				}
			}

			if !response.WasNotFound(existingCluster.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), clusterId)
			}

			dbId := databases.NewDatabaseID(subscriptionId, clusterId.ResourceGroupName, clusterId.RedisEnterpriseName, defaultDatabaseName)

			clusterParams := redisenterprise.Cluster{
				Location: location.Normalize(model.Location),
				Sku: redisenterprise.Sku{
					Name: redisenterprise.SkuName(model.SkuName),
				},
				Properties: &redisenterprise.ClusterCreateProperties{
					Encryption:          expandManagedRedisClusterCustomerManagedKey(model.CustomerManagedKey),
					MinimumTlsVersion:   pointer.To(redisenterprise.TlsVersionOnePointTwo),
					HighAvailability:    expandHighAvailability(model.HighAvailabilityEnabled),
					PublicNetworkAccess: redisenterprise.PublicNetworkAccess(model.PublicNetworkAccess),
				},
				Tags: pointer.To(model.Tags),
			}

			expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}
			clusterParams.Identity = expandedIdentity

			if err := clusterClient.CreateThenPoll(ctx, clusterId, clusterParams); err != nil {
				return fmt.Errorf("creating %s: %+v", clusterId, err)
			}

			metadata.SetID(clusterId)

			pollerType := custompollers.NewClusterStatePoller(clusterClient, clusterId)
			poller := pollers.NewPoller(pollerType, 15*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for `resourceState` to be `Running` for %s: %+v", clusterId, err)
			}

			if len(model.DefaultDatabase) == 1 {
				dbModel := model.DefaultDatabase[0]

				err := createDb(ctx, dbClient, dbId, dbModel)
				if err != nil {
					return fmt.Errorf("creating %s: %+v", dbId, err)
				}
			}

			return nil
		},
	}
}

func (r ManagedRedisResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			clusterClient := metadata.Client.ManagedRedis.Client
			dbClient := metadata.Client.ManagedRedis.DatabaseClient

			clusterId, err := redisenterprise.ParseRedisEnterpriseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			dbId := databases.NewDatabaseID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.RedisEnterpriseName, defaultDatabaseName)

			clusterResp, err := clusterClient.Get(ctx, *clusterId)
			if err != nil {
				if response.WasNotFound(clusterResp.HttpResponse) {
					return metadata.MarkAsGone(clusterId)
				}
				return fmt.Errorf("retrieving %s: %+v", clusterId, err)
			}

			state := ManagedRedisResourceModel{
				Name:              clusterId.RedisEnterpriseName,
				ResourceGroupName: clusterId.ResourceGroupName,
			}

			if model := clusterResp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.SkuName = string(model.Sku.Name)

				flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}

				state.Identity = pointer.From(flattenedIdentity)
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.CustomerManagedKey = flattenManagedRedisClusterCustomerManagedKey(props.Encryption)
					state.HighAvailabilityEnabled = strings.EqualFold(string(pointer.From(props.HighAvailability)), string(redisenterprise.HighAvailabilityEnabled))
					state.Hostname = pointer.From(props.HostName)
					state.PublicNetworkAccess = string(props.PublicNetworkAccess)
				}
			}

			dbResp, err := dbClient.Get(ctx, dbId)
			if err != nil {
				if !response.WasNotFound(dbResp.HttpResponse) {
					return fmt.Errorf("retrieving %s: %+v", dbId, err)
				}
			}

			if model := dbResp.Model; model != nil {
				if props := model.Properties; props != nil {
					defaultDb := DefaultDatabaseModel{
						AccessKeysAuthenticationEnabled: strings.EqualFold(pointer.FromEnum(props.AccessKeysAuthentication), string(databases.AccessKeysAuthenticationEnabled)),
						ClientProtocol:                  pointer.FromEnum(props.ClientProtocol),
						ClusteringPolicy:                pointer.FromEnum(props.ClusteringPolicy),
						EvictionPolicy:                  pointer.FromEnum(props.EvictionPolicy),
						GeoReplicationGroupName:         flattenGeoReplicationGroupName(props.GeoReplication),
						Module:                          flattenModules(props.Modules),
						Port:                            pointer.From(props.Port),
					}

					if defaultDb.AccessKeysAuthenticationEnabled {
						keysResp, err := dbClient.ListKeys(ctx, dbId)
						if err != nil {
							return fmt.Errorf("listing keys for %s: %+v", dbId, err)
						}

						if keysModel := keysResp.Model; keysModel != nil {
							defaultDb.PrimaryAccessKey = pointer.From(keysModel.PrimaryKey)
							defaultDb.SecondaryAccessKey = pointer.From(keysModel.SecondaryKey)
						}
					}

					state.DefaultDatabase = []DefaultDatabaseModel{defaultDb}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagedRedisResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			clusterClient := metadata.Client.ManagedRedis.Client
			dbClient := metadata.Client.ManagedRedis.DatabaseClient

			clusterId, err := redisenterprise.ParseRedisEnterpriseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			dbId := databases.NewDatabaseID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.RedisEnterpriseName, defaultDatabaseName)

			var state ManagedRedisResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existingCluster, err := clusterClient.Get(ctx, *clusterId)
			if err != nil {
				if response.WasNotFound(existingCluster.HttpResponse) {
					return metadata.MarkAsGone(clusterId)
				}

				return fmt.Errorf("retrieving existing %s: %+v", clusterId, err)
			}

			clusterParams := existingCluster.Model

			clusterUpdateRequired := false

			if metadata.ResourceData.HasChange("customer_managed_key") {
				clusterParams.Properties.Encryption = expandManagedRedisClusterCustomerManagedKey(state.CustomerManagedKey)
				clusterUpdateRequired = true
			}

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(state.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				clusterParams.Identity = expandedIdentity
				clusterUpdateRequired = true
			}

			if metadata.ResourceData.HasChange("public_network_access") {
				clusterParams.Properties.PublicNetworkAccess = redisenterprise.PublicNetworkAccess(state.PublicNetworkAccess)
				clusterUpdateRequired = true
			}

			if metadata.ResourceData.HasChange("tags") {
				clusterParams.Tags = pointer.To(state.Tags)
				clusterUpdateRequired = true
			}

			if clusterUpdateRequired {
				// Despite the method name, Create uses PUT (create-or-update behaviour), which is preferred to Update (PATCH)
				// to simplify 'omitempty' / empty values handling on expand functions
				if err := clusterClient.CreateThenPoll(ctx, *clusterId, *clusterParams); err != nil {
					return fmt.Errorf("creating cluster %s: %+v", clusterId, err)
				}

				pollerType := custompollers.NewClusterStatePoller(clusterClient, *clusterId)
				poller := pollers.NewPoller(pollerType, 15*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
				if err := poller.PollUntilDone(ctx); err != nil {
					return fmt.Errorf("waiting for `resourceState` to be `Running` for %s: %+v", clusterId, err)
				}
			}

			if metadata.ResourceData.HasChange("default_database") {
				old, new := metadata.ResourceData.GetChange("default_database")
				switch {
				case dbLen(old) == 1 && dbLen(new) == 0:
					if err := dbClient.DeleteThenPoll(ctx, dbId); err != nil {
						return fmt.Errorf("deleting database %s: %+v", dbId, err)
					}
				case dbLen(old) == 0 && dbLen(new) == 1:
					if err := createDb(ctx, dbClient, dbId, state.DefaultDatabase[0]); err != nil {
						return fmt.Errorf("creating %s: %+v", dbId, err)
					}
				default:
					if metadata.ResourceData.HasChanges(
						"default_database.0.clustering_policy",
						"default_database.0.geo_replication_group_name",
						"default_database.0.module",
					) {
						log.Printf("[INFO] re-creating database %s to apply updates to immutable properties, data will be lost and Managed Redis will be unavailable during this operation", dbId)

						if err := dbClient.DeleteThenPoll(ctx, dbId); err != nil {
							return fmt.Errorf("deleting database %s for re-creation: %+v", dbId, err)
						}
						if err := createDb(ctx, dbClient, dbId, state.DefaultDatabase[0]); err != nil {
							return fmt.Errorf("re-creating %s: %+v", dbId, err)
						}
					} else if metadata.ResourceData.HasChanges(
						"default_database.0.access_keys_authentication_enabled",
						"default_database.0.client_protocol",
						"default_database.0.eviction_policy",
					) {
						existingDb, err := dbClient.Get(ctx, dbId)
						if err != nil {
							return fmt.Errorf("retrieving existing %s: %+v", dbId, err)
						}

						dbParams := existingDb.Model

						if dbParams == nil {
							return fmt.Errorf("retrieving existing %s: `model` was nil", dbId)
						}
						if dbParams.Properties == nil {
							return fmt.Errorf("retrieving existing %s: `properties` was nil", dbId)
						}

						if metadata.ResourceData.HasChange("default_database.0.access_keys_authentication_enabled") {
							dbParams.Properties.AccessKeysAuthentication = expandAccessKeysAuth(state.DefaultDatabase[0].AccessKeysAuthenticationEnabled)
						}
						if metadata.ResourceData.HasChange("default_database.0.client_protocol") {
							dbParams.Properties.ClientProtocol = pointer.ToEnum[databases.Protocol](state.DefaultDatabase[0].ClientProtocol)
						}
						if metadata.ResourceData.HasChange("default_database.0.eviction_policy") {
							dbParams.Properties.EvictionPolicy = pointer.ToEnum[databases.EvictionPolicy](state.DefaultDatabase[0].EvictionPolicy)
						}

						// Despite the method name, Create uses PUT (create-or-update behaviour), which is preferred to Update (PATCH)
						// to simplify 'omitempty' / empty values handling on expand functions
						if err := dbClient.CreateThenPoll(ctx, dbId, *dbParams); err != nil {
							return fmt.Errorf("updating %s: %+v", dbId, err)
						}

						pollerType := custompollers.NewDBStatePoller(dbClient, dbId)
						poller := pollers.NewPoller(pollerType, 15*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
						if err := poller.PollUntilDone(ctx); err != nil {
							return fmt.Errorf("waiting for `resourceState` to be `Running` for %s: %+v", dbId, err)
						}
					}
				}
			}

			return nil
		},
	}
}

func (r ManagedRedisResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			clusterClient := metadata.Client.ManagedRedis.Client
			dbClient := metadata.Client.ManagedRedis.DatabaseClient

			clusterId, err := redisenterprise.ParseRedisEnterpriseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			dbId := databases.NewDatabaseID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.RedisEnterpriseName, defaultDatabaseName)

			if err := dbClient.DeleteThenPoll(ctx, dbId); err != nil {
				return fmt.Errorf("deleting %s: %+v", dbId, err)
			}

			if err := clusterClient.DeleteThenPoll(ctx, *clusterId); err != nil {
				return fmt.Errorf("deleting %s: %+v", clusterId, err)
			}

			return nil
		},
	}
}

func (r ManagedRedisResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}

			var model ManagedRedisResourceModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return err
			}

			if len(model.DefaultDatabase) > 0 {
				dbModel := model.DefaultDatabase[0]

				if geoReplicationEnabled := dbModel.GeoReplicationGroupName != ""; geoReplicationEnabled {
					if !slices.Contains(validate.SKUsSupportingGeoReplication(), model.SkuName) {
						return fmt.Errorf("SKU %q does not support geo-replication, only following SKUs are supported: %s", model.SkuName, strings.Join(validate.SKUsSupportingGeoReplication(), ", "))
					}

					for _, module := range dbModel.Module {
						if module.Name != "" && !slices.Contains(validate.DatabaseModulesSupportingGeoReplication(), module.Name) {
							return fmt.Errorf("invalid module %q, only following modules are supported when `geo_replication_group_name` is not empty: %s", module.Name, strings.Join(validate.DatabaseModulesSupportingGeoReplication(), ", "))
						}
					}
				}

				if dbModel.EvictionPolicy != "" {
					for _, module := range dbModel.Module {
						if module.Name != "" && module.Name == "RediSearch" {
							if dbModel.EvictionPolicy != string(redisenterprise.EvictionPolicyNoEviction) {
								return fmt.Errorf("invalid eviction_policy %q, when using RediSearch module, eviction_policy must be set to NoEviction", dbModel.EvictionPolicy)
							}

							if dbModel.ClusteringPolicy != string(redisenterprise.ClusteringPolicyEnterpriseCluster) {
								return fmt.Errorf("invalid clustering_policy %q, when using RediSearch module, clustering_policy must be set to EnterpriseCluster", dbModel.ClusteringPolicy)
							}
						}
					}
				}
			}

			return nil
		},
	}
}

func createDb(ctx context.Context, dbClient *databases.DatabasesClient, dbId databases.DatabaseId, dbModel DefaultDatabaseModel) error {
	dbParams := databases.Database{
		Properties: &databases.DatabaseCreateProperties{
			AccessKeysAuthentication: expandAccessKeysAuth(dbModel.AccessKeysAuthenticationEnabled),
			ClientProtocol:           pointer.To(databases.Protocol(dbModel.ClientProtocol)),
			ClusteringPolicy:         pointer.To(databases.ClusteringPolicy(dbModel.ClusteringPolicy)),
			EvictionPolicy:           pointer.To(databases.EvictionPolicy(dbModel.EvictionPolicy)),
			GeoReplication:           expandGeoReplication(dbModel.GeoReplicationGroupName, dbId.ID()),
			Modules:                  expandModules(dbModel.Module),
		},
	}

	if err := dbClient.CreateThenPoll(ctx, dbId, dbParams); err != nil {
		return fmt.Errorf("creating database %s: %+v", dbId, err)
	}

	pollerType := custompollers.NewDBStatePoller(dbClient, dbId)
	poller := pollers.NewPoller(pollerType, 15*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for `resourceState` to be `Running` for %s: %+v", dbId, err)
	}
	return nil
}

func expandManagedRedisClusterCustomerManagedKey(input []CustomerManagedKeyModel) *redisenterprise.ClusterCommonPropertiesEncryption {
	if len(input) == 0 {
		return &redisenterprise.ClusterCommonPropertiesEncryption{}
	}

	cmk := input[0]

	return &redisenterprise.ClusterCommonPropertiesEncryption{
		CustomerManagedKeyEncryption: &redisenterprise.ClusterCommonPropertiesEncryptionCustomerManagedKeyEncryption{
			KeyEncryptionKeyURL: pointer.To(cmk.KeyVaultKeyId),
			KeyEncryptionKeyIdentity: &redisenterprise.ClusterCommonPropertiesEncryptionCustomerManagedKeyEncryptionKeyEncryptionKeyIdentity{
				IdentityType:                   pointer.To(redisenterprise.CmkIdentityTypeUserAssignedIdentity),
				UserAssignedIdentityResourceId: pointer.To(cmk.UserAssignedIdentityId),
			},
		},
	}
}

func expandHighAvailability(enabled bool) *redisenterprise.HighAvailability {
	if enabled {
		return pointer.To(redisenterprise.HighAvailabilityEnabled)
	}

	return pointer.To(redisenterprise.HighAvailabilityDisabled)
}

func expandAccessKeysAuth(enabled bool) *databases.AccessKeysAuthentication {
	if enabled {
		return pointer.To(databases.AccessKeysAuthenticationEnabled)
	}

	return pointer.To(databases.AccessKeysAuthenticationDisabled)
}

func expandGeoReplication(input string, id string) *databases.DatabaseCommonPropertiesGeoReplication {
	if input == "" {
		return nil
	}

	return &databases.DatabaseCommonPropertiesGeoReplication{
		GroupNickname: pointer.To(input),
		LinkedDatabases: &[]databases.LinkedDatabase{
			{
				Id: pointer.To(id),
			},
		},
	}
}

func flattenGeoReplicationGroupName(input *databases.DatabaseCommonPropertiesGeoReplication) string {
	if input == nil || input.GroupNickname == nil {
		return ""
	}
	return pointer.From(input.GroupNickname)
}

func expandModules(input []ModuleModel) *[]databases.Module {
	results := make([]databases.Module, 0, len(input))
	for _, module := range input {
		results = append(results, databases.Module{
			Name: module.Name,
			Args: pointer.To(module.Args),
		})
	}
	return &results
}

func flattenModules(input *[]databases.Module) []ModuleModel {
	results := make([]ModuleModel, 0)
	if input == nil {
		return results
	}

	for _, module := range *input {
		results = append(results, ModuleModel{
			Name:    module.Name,
			Args:    pointer.From(module.Args),
			Version: pointer.From(module.Version),
		})
	}
	return results
}

func flattenManagedRedisClusterCustomerManagedKey(input *redisenterprise.ClusterCommonPropertiesEncryption) []CustomerManagedKeyModel {
	if input == nil || input.CustomerManagedKeyEncryption == nil {
		return []CustomerManagedKeyModel{}
	}

	cmkEncryption := input.CustomerManagedKeyEncryption
	uaiResourceId := ""
	if cmkEncryption.KeyEncryptionKeyIdentity != nil {
		uaiResourceId = pointer.From(cmkEncryption.KeyEncryptionKeyIdentity.UserAssignedIdentityResourceId)
	}

	return []CustomerManagedKeyModel{
		{
			KeyVaultKeyId:          pointer.From(cmkEncryption.KeyEncryptionKeyURL),
			UserAssignedIdentityId: uaiResourceId,
		},
	}
}

func dbLen(v interface{}) int {
	if s, ok := v.([]interface{}); ok {
		return len(s)
	}
	return 0
}
