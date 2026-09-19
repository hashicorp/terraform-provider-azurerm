// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
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

			state := ResourceModel{
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
					databaseId, err := redisenterprise.ParseDatabaseID(pointer.From(model.Id))
					if err != nil {
						return fmt.Errorf("parsing Managed Redis Database ID %q: %+v", pointer.From(model.Id), err)
					}

					defaultDb := DefaultDatabaseModel{
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

func flattenGeoReplicationGroupName(input *databases.DatabasePropertiesGeoReplication) string {
	if input == nil || input.GroupNickname == nil {
		return ""
	}
	return pointer.From(input.GroupNickname)
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

func flattenManagedRedisClusterCustomerManagedKey(input *redisenterprise.ClusterPropertiesEncryption) []CustomerManagedKeyModel {
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

func flattenPersistenceAOF(input *databases.Persistence) string {
	if input == nil {
		return ""
	}

	if pointer.From(input.AofEnabled) {
		return pointer.FromEnum(input.AofFrequency)
	}

	return ""
}

func flattenPersistenceRDB(input *databases.Persistence) string {
	if input == nil {
		return ""
	}

	if pointer.From(input.RdbEnabled) {
		return pointer.FromEnum(input.RdbFrequency)
	}

	return ""
}
