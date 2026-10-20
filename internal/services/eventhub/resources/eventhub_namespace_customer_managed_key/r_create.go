// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace_customer_managed_key

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/namespaces"
)

func expandEventHubNamespaceKeyVaultKeyIds(input []any) (*[]namespaces.KeyVaultProperties, error) {
	if len(input) == 0 {
		return nil, nil
	}

	results := make([]namespaces.KeyVaultProperties, 0)

	for _, item := range input {
		keyId, err := keyvault.ParseNestedItemID(item.(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
		if err != nil {
			return nil, err
		}

		results = append(results, namespaces.KeyVaultProperties{
			KeyName:     pointer.To(keyId.Name),
			KeyVaultUri: pointer.To(keyId.KeyVaultBaseURL),
			KeyVersion:  pointer.To(keyId.Version),
		})
	}

	return &results, nil
}
