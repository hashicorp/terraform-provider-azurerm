// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package managed_redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataSource struct{}

var _ sdk.DataSource = DataSource{}

type ManagedRedisDataSourceModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`

	CustomerManagedKey      []CustomerManagedKeyModel                  `tfschema:"customer_managed_key"`
	DefaultDatabase         []DefaultDatabaseDataSourceModel           `tfschema:"default_database"`
	HighAvailabilityEnabled bool                                       `tfschema:"high_availability_enabled"`
	Hostname                string                                     `tfschema:"hostname"`
	Identity                []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Location                string                                     `tfschema:"location"`
	PublicNetworkAccess     string                                     `tfschema:"public_network_access"`
	SkuName                 string                                     `tfschema:"sku_name"`
	Tags                    map[string]string                          `tfschema:"tags"`
}

type DefaultDatabaseDataSourceModel struct {
	AccessKeysAuthenticationEnabled          bool          `tfschema:"access_keys_authentication_enabled"`
	ClientProtocol                           string        `tfschema:"client_protocol"`
	ClusteringPolicy                         string        `tfschema:"clustering_policy"`
	EvictionPolicy                           string        `tfschema:"eviction_policy"`
	GeoReplicationGroupName                  string        `tfschema:"geo_replication_group_name"`
	GeoReplicationLinkedDatabaseIds          []string      `tfschema:"geo_replication_linked_database_ids"`
	ID                                       string        `tfschema:"id"`
	Module                                   []ModuleModel `tfschema:"module"`
	PersistenceAppendOnlyFileBackupFrequency string        `tfschema:"persistence_append_only_file_backup_frequency"`
	PersistenceRedisDatabaseBackupFrequency  string        `tfschema:"persistence_redis_database_backup_frequency"`
	Port                                     int64         `tfschema:"port"`
	PrimaryAccessKey                         string        `tfschema:"primary_access_key"`
	SecondaryAccessKey                       string        `tfschema:"secondary_access_key"`
}

func (r DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Managed Redis instance.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ManagedRedisClusterName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"customer_managed_key": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Description: "The ID of the Key Vault Key used for encryption.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"user_assigned_identity_id": {
						Description: "The ID of the User Assigned Identity that has access to the Key Vault Key.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"default_database": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Description: "The ID of the Managed Redis Database Instance.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"access_keys_authentication_enabled": {
						Description: "Whether access key authentication is enabled for the database.",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},

					"client_protocol": {
						Description: "The client protocol used by the database (either `Encrypted` or `Plaintext`).",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"clustering_policy": {
						Description: "The clustering policy used by the database.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"eviction_policy": {
						Description: "The Redis eviction policy used by the database.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"geo_replication_group_name": {
						Description: "The name of the geo-replication group.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"geo_replication_linked_database_ids": {
						Description: "A list of linked database IDs for geo-replication.",
						Type:        pluginsdk.TypeSet,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"module": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description: "The name of the Redis module.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},

								"args": {
									Description: "The configuration options for the module.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},

								"version": {
									Description: "The version of the module.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
							},
						},
					},

					"persistence_append_only_file_backup_frequency": {
						Description: "The frequency of Append Only File (AOF) backups.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"persistence_redis_database_backup_frequency": {
						Description: "The frequency of Redis Database (RDB) backups.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"port": {
						Description: "The TCP port of the database endpoint.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},

					"primary_access_key": {
						Description: "The Primary Access Key for the Managed Redis Database instance.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},

					"secondary_access_key": {
						Description: "The Secondary Access Key for the Managed Redis Database instance.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
		},

		"high_availability_enabled": {
			Description: "Whether high availability is enabled for the Managed Redis instance.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"hostname": {
			Description: "The DNS hostname of the Managed Redis instance.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

		"location": commonschema.LocationComputed(),

		"public_network_access": {
			Description: "The public network access setting for the Managed Redis instance.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"sku_name": {
			Description: "The SKU name of the Managed Redis instance.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r DataSource) ModelObject() any {
	return &ManagedRedisDataSourceModel{}
}

func (r DataSource) ResourceType() string {
	return "azurerm_managed_redis"
}

func (r DataSource) Description() string {
	return "Use this data source to access information about an existing Managed Redis instance."
}

func (r DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			clusterClient := metadata.Client.ManagedRedis.Client
			dbClient := metadata.Client.ManagedRedis.DatabaseClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state ManagedRedisDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			clusterId := redisenterprise.NewRedisEnterpriseID(subscriptionId, state.ResourceGroupName, state.Name)
			dbId := databases.NewDatabaseID(subscriptionId, clusterId.ResourceGroupName, clusterId.RedisEnterpriseName, defaultDatabaseName)

			clusterResp, err := clusterClient.Get(ctx, clusterId)
			if err != nil {
				if response.WasNotFound(clusterResp.HttpResponse) {
					return fmt.Errorf("%s was not found", clusterId)
				}
				return fmt.Errorf("retrieving %s: %+v", clusterId, err)
			}

			dbResp, err := dbClient.Get(ctx, dbId)
			if err != nil {
				if !response.WasNotFound(dbResp.HttpResponse) {
					return fmt.Errorf("retrieving %s: %+v", dbId, err)
				}
			}

			metadata.SetID(clusterId)

			if model := clusterResp.Model; model != nil {
				flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}

				state.Identity = pointer.From(flattenedIdentity)
				state.Location = location.Normalize(model.Location)
				state.SkuName = string(model.Sku.Name)
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.CustomerManagedKey = flattenManagedRedisClusterCustomerManagedKey(props.Encryption)
					state.HighAvailabilityEnabled = strings.EqualFold(pointer.FromEnum(props.HighAvailability), string(redisenterprise.HighAvailabilityEnabled))
					state.Hostname = pointer.From(props.HostName)
					state.PublicNetworkAccess = string(props.PublicNetworkAccess)
				}
			}

			if model := dbResp.Model; model != nil {
				if props := model.Properties; props != nil {
					databaseId, err := redisenterprise.ParseDatabaseID(pointer.From(model.Id))
					if err != nil {
						return err
					}

					defaultDb := DefaultDatabaseDataSourceModel{
						AccessKeysAuthenticationEnabled:          strings.EqualFold(pointer.FromEnum(props.AccessKeysAuthentication), string(databases.AccessKeysAuthenticationEnabled)),
						ClientProtocol:                           pointer.FromEnum(props.ClientProtocol),
						ClusteringPolicy:                         pointer.FromEnum(props.ClusteringPolicy),
						EvictionPolicy:                           pointer.FromEnum(props.EvictionPolicy),
						GeoReplicationGroupName:                  flattenGeoReplicationGroupName(props.GeoReplication),
						ID:                                       databaseId.ID(),
						Module:                                   flattenModules(props.Modules),
						PersistenceAppendOnlyFileBackupFrequency: flattenPersistenceAOF(props.Persistence),
						PersistenceRedisDatabaseBackupFrequency:  flattenPersistenceRDB(props.Persistence),
						Port:                                     pointer.From(props.Port),
					}

					if props.GeoReplication != nil {
						defaultDb.GeoReplicationLinkedDatabaseIds = helpers.FlattenLinkedDatabases(props.GeoReplication.LinkedDatabases)
					}

					if defaultDb.AccessKeysAuthenticationEnabled {
						keysResp, err := dbClient.ListKeys(ctx, dbId)
						if err != nil {
							return fmt.Errorf("listing keys for %s: %+v", dbId, err)
						}
						if keysModel := keysResp.Model; keysModel != nil {
							defaultDb.PrimaryAccessKey = pointer.From(keysResp.Model.PrimaryKey)
							defaultDb.SecondaryAccessKey = pointer.From(keysResp.Model.SecondaryKey)
						}
					}

					state.DefaultDatabase = []DefaultDatabaseDataSourceModel{defaultDb}
				}
			}

			return metadata.Encode(&state)
		},
	}
}
