package azure

import (
	keyVaultParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
)

type KeyVaultChildID = keyVaultParse.NestedItemId

func ParseKeyVaultChildID(id string) (*KeyVaultChildID, error) {
	return keyVaultParse.ParseNestedItemID(id)
}

func ParseKeyVaultChildIDVersionOptional(id string) (*KeyVaultChildID, error) {
	return keyVaultParse.ParseOptionallyVersionedNestedItemID(id)
}

func ValidateKeyVaultChildName(v interface{}, k string) (warnings []string, errors []error) {
	return keyVaultValidate.NestedItemName(v, k)
}

func ValidateKeyVaultChildId(i interface{}, k string) (warnings []string, errors []error) {
	return keyVaultValidate.NestedItemId(i, k)
}
