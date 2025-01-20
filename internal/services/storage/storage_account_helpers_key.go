// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	managedHsmParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
)

type accountKeyDetails struct {
	keyVaultBaseUrl  string
	keyVaultKeyUri   string
	managedHsmKeyUri string
	keyName          string
	keyVersion       string
}

// nolint unparam // keyVaultApi may be used in future
func flattenCustomerManagedKey(input *storageaccounts.KeyVaultProperties, keyVaultApi, managedHsmApi environments.Api) accountKeyDetails {
	output := accountKeyDetails{
		keyVaultBaseUrl:  "",
		keyVaultKeyUri:   "",
		managedHsmKeyUri: "",
		keyName:          "",
		keyVersion:       "",
	}

	if input == nil || input.Keyvaulturi == nil || input.Keyname == nil {
		return output
	}

	// Whilst this says Key Vault it contains either a Key Vault or Managed HSM Key ID
	baseUri := pointer.From(input.Keyvaulturi)
	output.keyName = pointer.From(input.Keyname)
	output.keyVersion = pointer.From(input.Keyversion)
	itemId := fmt.Sprintf("%s/keys/%s", strings.TrimSuffix(baseUri, "/"), output.keyName)

	// This either has no version (i.e. use latest)
	if output.keyVersion == "" {
		parsedKeyVaultId, _ := keyVaultParse.ParseOptionallyVersionedNestedItemID(itemId)
		if parsedKeyVaultId != nil {
			output.keyVaultBaseUrl = baseUri
			output.keyVaultKeyUri = parsedKeyVaultId.ID()
			return output
		}

		if domainSuffix, ok := managedHsmApi.DomainSuffix(); ok {
			if parsedManagedHsmId, _ := managedHsmParse.ManagedHSMDataPlaneVersionlessKeyID(itemId, domainSuffix); parsedManagedHsmId != nil {
				output.managedHsmKeyUri = parsedManagedHsmId.ID()
				return output
			}
		}
	}

	// or the key is for a specific version of a key
	if output.keyVersion != "" {
		itemId = fmt.Sprintf("%s/%s", itemId, output.keyVersion)

		parsedKeyVaultId, _ := keyVaultParse.ParseNestedItemID(itemId)
		if parsedKeyVaultId != nil {
			output.keyVaultBaseUrl = baseUri
			output.keyVaultKeyUri = parsedKeyVaultId.ID()
			return output
		}

		if domainSuffix, ok := managedHsmApi.DomainSuffix(); ok {
			if parsedManagedHsmId, _ := managedHsmParse.ManagedHSMDataPlaneVersionedKeyID(itemId, domainSuffix); parsedManagedHsmId != nil {
				output.managedHsmKeyUri = parsedManagedHsmId.ID()
				return output
			}
		}
	}

	return output
}
