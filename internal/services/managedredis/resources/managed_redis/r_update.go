// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/custompollers"
)

func (r Resource) Update() sdk.ResourceFunc {
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

			if metadata.ResourceData.HasChange("sku_name") {
				clusterParams.Sku.Name = redisenterprise.SkuName(state.SkuName)
				clusterUpdateRequired = true
			}

			if metadata.ResourceData.HasChange("tags") {
				clusterParams.Tags = new(state.Tags)
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
					// lintignore:R019 // deliberate subsets: the first branch lists fields requiring database re-creation, the second lists fields updatable in place
					if metadata.ResourceData.HasChanges(
						"default_database.0.clustering_policy",
						"default_database.0.geo_replication_group_name",
						"default_database.0.module",
					) {
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
						"default_database.0.persistence_append_only_file_backup_frequency",
						"default_database.0.persistence_redis_database_backup_frequency",
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
						if metadata.ResourceData.HasChanges(
							"default_database.0.persistence_append_only_file_backup_frequency",
							"default_database.0.persistence_redis_database_backup_frequency",
						) {
							dbParams.Properties.Persistence = expandPersistence(
								state.DefaultDatabase[0].PersistenceAppendOnlyFileBackupFrequency,
								state.DefaultDatabase[0].PersistenceRedisDatabaseBackupFrequency,
							)
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
