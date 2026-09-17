// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_redis

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
)

func expandCreateForManagedRedis(model ResourceModel) (redisenterprise.Cluster, error) {
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
			KeyEncryptionKeyURL: pointer.To(cmk.KeyVaultKeyId),
			KeyEncryptionKeyIdentity: &redisenterprise.ClusterPropertiesEncryptionCustomerManagedKeyEncryptionKeyEncryptionKeyIdentity{
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

func expandGeoReplication(input string, id string) *databases.DatabasePropertiesGeoReplication {
	if input == "" {
		return nil
	}

	return &databases.DatabasePropertiesGeoReplication{
		GroupNickname: pointer.To(input),
		LinkedDatabases: &[]databases.LinkedDatabase{
			{
				Id: pointer.To(id),
			},
		},
	}
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

func expandPersistence(aofBackupFreq string, rdbBackupFreq string) *databases.Persistence {
	switch {
	case aofBackupFreq != "":
		return &databases.Persistence{
			AofEnabled:   pointer.To(true),
			AofFrequency: pointer.ToEnum[databases.AofFrequency](aofBackupFreq),
		}
	case rdbBackupFreq != "":
		return &databases.Persistence{
			RdbEnabled:   pointer.To(true),
			RdbFrequency: pointer.ToEnum[databases.RdbFrequency](rdbBackupFreq),
		}
	default:
		return &databases.Persistence{}
	}
}
