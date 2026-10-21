// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package suppress

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func DiffSuppressIgnoreKeyVaultKeyVersion(k, old, new string, _ *pluginsdk.ResourceData) bool {
	// TODO: deprecate this method in the future, `ignore_changes` should be used instead
	oldKey, err := keyvault.ParseNestedItemID(old, keyvault.VersionTypeAny, keyvault.NestedItemTypeAny)
	if err != nil {
		return false
	}
	newKey, err := keyvault.ParseNestedItemID(new, keyvault.VersionTypeAny, keyvault.NestedItemTypeAny)
	if err != nil {
		return false
	}

	return (oldKey.KeyVaultBaseURL == newKey.KeyVaultBaseURL) && (oldKey.Name == newKey.Name)
}
