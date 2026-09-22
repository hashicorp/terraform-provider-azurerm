// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/custompollers"
)

func expandCreateForManagedRedis(model ManagedRedisResourceModel) (redisenterprise.Cluster, error) {
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
		Tags: new(model.Tags),
	}

	expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
	if err != nil {
		return clusterParams, fmt.Errorf("expanding `identity`: %+v", err)
	}
	clusterParams.Identity = expandedIdentity

	return clusterParams, nil
}

func expandManagedRedisClusterCustomerManagedKey(input []CustomerManagedKeyModel) *redisenterprise.ClusterPropertiesEncryption {
	if len(input) == 0 {
		return &redisenterprise.ClusterPropertiesEncryption{}
	}

	cmk := input[0]

	return &redisenterprise.ClusterPropertiesEncryption{
		CustomerManagedKeyEncryption: &redisenterprise.ClusterPropertiesEncryptionCustomerManagedKeyEncryption{
			KeyEncryptionKeyURL: new(cmk.KeyVaultKeyId),
			KeyEncryptionKeyIdentity: &redisenterprise.ClusterPropertiesEncryptionCustomerManagedKeyEncryptionKeyEncryptionKeyIdentity{
				IdentityType:                   pointer.To(redisenterprise.CmkIdentityTypeUserAssignedIdentity),
				UserAssignedIdentityResourceId: new(cmk.UserAssignedIdentityId),
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

func expandGeoReplication(input string, id string) *databases.DatabasePropertiesGeoReplication {
	if input == "" {
		return nil
	}

	return &databases.DatabasePropertiesGeoReplication{
		GroupNickname: new(input),
		LinkedDatabases: &[]databases.LinkedDatabase{
			{
				Id: new(id),
			},
		},
	}
}

func expandModules(input []ModuleModel) *[]databases.Module {
	results := make([]databases.Module, 0, len(input))
	for _, module := range input {
		results = append(results, databases.Module{
			Name: module.Name,
			Args: new(module.Args),
		})
	}
	return &results
}

func expandPersistence(aofBackupFreq string, rdbBackupFreq string) *databases.Persistence {
	switch {
	case aofBackupFreq != "":
		return &databases.Persistence{
			AofEnabled:   new(true),
			AofFrequency: pointer.ToEnum[databases.AofFrequency](aofBackupFreq),
		}
	case rdbBackupFreq != "":
		return &databases.Persistence{
			RdbEnabled:   new(true),
			RdbFrequency: pointer.ToEnum[databases.RdbFrequency](rdbBackupFreq),
		}
	default:
		return &databases.Persistence{}
	}
}

func createDb(ctx context.Context, dbClient *databases.DatabasesClient, dbId databases.DatabaseId, dbModel DefaultDatabaseModel) error {
	dbParams := databases.Database{
		Properties: &databases.DatabaseCreateProperties{
			AccessKeysAuthentication: expandAccessKeysAuth(dbModel.AccessKeysAuthenticationEnabled),
			ClientProtocol:           pointer.ToEnum[databases.Protocol](dbModel.ClientProtocol),
			ClusteringPolicy:         pointer.ToEnum[databases.ClusteringPolicy](dbModel.ClusteringPolicy),
			EvictionPolicy:           pointer.ToEnum[databases.EvictionPolicy](dbModel.EvictionPolicy),
			GeoReplication:           expandGeoReplication(dbModel.GeoReplicationGroupName, dbId.ID()),
			Modules:                  expandModules(dbModel.Module),
			Persistence:              expandPersistence(dbModel.PersistenceAppendOnlyFileBackupFrequency, dbModel.PersistenceRedisDatabaseBackupFrequency),
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
