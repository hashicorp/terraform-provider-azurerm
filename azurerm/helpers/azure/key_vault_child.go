package azure

import (
	keyVaultParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
)

type KeyVaultChildID = keyVaultParse.NestedItemId

func ParseKeyVaultChildID(id string) (*KeyVaultChildID, error) {
	return keyVaultParse.ParseNestedItemID(id)
}
