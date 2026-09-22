// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_redis

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/preflight"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
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
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

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

					"persistence_append_only_file_backup_frequency": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ConflictsWith: []string{
							"default_database.0.geo_replication_group_name",
							"default_database.0.persistence_redis_database_backup_frequency",
						},

						ValidateFunc: validation.StringInSlice(validate.PossibleValuesForAofFrequency(), false),
					},

					"persistence_redis_database_backup_frequency": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ConflictsWith: []string{
							"default_database.0.geo_replication_group_name",
							"default_database.0.persistence_append_only_file_backup_frequency",
						},
						ValidateFunc: validation.StringInSlice(redisenterprise.PossibleValuesForRdbFrequency(), false),
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

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"hostname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
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

			if metadata.Client.Features.EnhancedValidation.PreflightEnabled {
				// Only perform preflight validation if there are changes. This avoids validation failures and
				// additional API calls for resources that are unchanged between plan invocations
				if len(metadata.ResourceDiff.GetChangedKeysPrefix("")) > 0 || metadata.ResourceDiff.Id() == "" {
					req, err := expandCreateForManagedRedis(model)
					if err != nil {
						return err
					}

					resId := redisenterprise.NewRedisEnterpriseID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)
					preflightValidate, err := preflight.NewValidationRequestWithTypeOverride(pointer.To(model.Location), pointer.To(resId), "redis", "2025-07-01", req)
					if err != nil {
						return fmt.Errorf("constructing preflight validation request: %w", err)
					}

					if err = preflightValidate.ValidateResource(ctx, metadata); err != nil {
						return err
					}
				}
			}

			if metadata.ResourceDiff.Id() == "" && len(model.DefaultDatabase) == 0 {
				return fmt.Errorf("`default_database` must be provided when creating a new resource")
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

			if metadata.ResourceDiff.HasChanges("sku_name") {
				if resId := metadata.ResourceDiff.Id(); resId != "" {
					clusterId, err := redisenterprise.ParseRedisEnterpriseID(resId)
					if err != nil {
						return err
					}
					clusterClient := metadata.Client.ManagedRedis.Client
					if !isSkuAllowedForScaling(ctx, clusterClient, clusterId, model.SkuName) {
						metadata.ResourceDiff.ForceNew("sku_name")
					}
				}
			}

			return nil
		},
	}
}
