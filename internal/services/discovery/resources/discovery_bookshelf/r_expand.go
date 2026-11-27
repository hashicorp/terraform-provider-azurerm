// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package discovery_bookshelf

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/bookshelves"
)

func expandDiscoveryBookshelfCustomerManagedKey(input []DiscoveryBookshelfCustomerManagedKey, identityClientId string) (*bookshelves.BookshelfKeyVaultProperties, error) {
	if len(input) == 0 {
		return nil, nil
	}

	keyId, err := keyvault.ParseNestedItemID(input[0].KeyVaultKeyId, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
	if err != nil {
		return nil, err
	}

	return &bookshelves.BookshelfKeyVaultProperties{
		IdentityClientId: identityClientId,
		KeyName:          keyId.Name,
		KeyVaultUri:      keyId.KeyVaultBaseURL,
		KeyVersion:       pointer.To(keyId.Version),
	}, nil
}

func expandDiscoveryBookshelfWorkloadIdentities(input []string) *map[string]bookshelves.UserAssignedIdentity {
	if len(input) == 0 {
		return nil
	}

	output := make(map[string]bookshelves.UserAssignedIdentity, len(input))
	for _, id := range input {
		output[id] = bookshelves.UserAssignedIdentity{}
	}

	return &output
}
