// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_redis

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
)

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
