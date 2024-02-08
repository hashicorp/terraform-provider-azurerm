// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package suppress

import (
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func DiffSuppressIgnoreKeyVaultKeyVersion(k, old, new string, d *pluginsdk.ResourceData) bool {
	// TODO: deprecate this method in the future, `ignore_changes` should be used instead
	oldKey, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(old)
	if err != nil {
		return false
	}
	newKey, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(new)
	if err != nil {
		return false
	}

	return (oldKey.KeyVaultBaseUrl == newKey.KeyVaultBaseUrl) && (oldKey.Name == newKey.Name)
}
