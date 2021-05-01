package suppress

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	keyVaultParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
)

func DiffSuppressIgnoreKeyVaultKeyVersion(k, old, new string, d *schema.ResourceData) bool {
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
