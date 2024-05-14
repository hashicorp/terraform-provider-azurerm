package storage

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	managedHsmParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
)

type accountKeyDetails struct {
	keyVaultKeyId string
	managedHsmId  string
}

func flattenCustomerManagedKey(input *storageaccounts.KeyVaultProperties, managedHsmApi environments.Api) accountKeyDetails {
	if input == nil {
		return accountKeyDetails{
			keyVaultKeyId: "",
			managedHsmId:  "",
		}
	}

	// Whilst this says Key Vault it contains either a Key Vault or Managed HSM Key ID
	baseUri := pointer.From(input.Keyvaulturi)
	keyName := pointer.From(input.Keyname)
	keyVersion := pointer.From(input.Keyversion)
	itemId := fmt.Sprintf("%s/keys/%s", baseUri, keyName)

	// This either has no version (i.e. use latest)
	if keyVersion == "" {
		parsedKeyVaultId, _ := keyVaultParse.ParseOptionallyVersionedNestedItemID(itemId)
		if parsedKeyVaultId != nil {
			return accountKeyDetails{
				keyVaultKeyId: parsedKeyVaultId.ID(),
				managedHsmId:  "",
			}
		}

		if domainSuffix, ok := managedHsmApi.DomainSuffix(); ok {
			if parsedManagedHsmId, _ := managedHsmParse.ManagedHSMDataPlaneVersionlessKeyID(itemId, domainSuffix); parsedManagedHsmId != nil {
				return accountKeyDetails{
					keyVaultKeyId: "",
					managedHsmId:  parsedManagedHsmId.ID(),
				}
			}
		}
	}

	// or the key is for a specific version of a key
	if keyVersion != "" {
		itemId = fmt.Sprintf("%s/%s", itemId, keyVersion)

		parsedKeyVaultId, _ := keyVaultParse.ParseNestedItemID(itemId)
		if parsedKeyVaultId != nil {
			return accountKeyDetails{
				keyVaultKeyId: parsedKeyVaultId.ID(),
				managedHsmId:  "",
			}
		}

		if domainSuffix, ok := managedHsmApi.DomainSuffix(); ok {
			if parsedManagedHsmId, _ := managedHsmParse.ManagedHSMDataPlaneVersionedKeyID(itemId, domainSuffix); parsedManagedHsmId != nil {
				return accountKeyDetails{
					keyVaultKeyId: "",
					managedHsmId:  parsedManagedHsmId.ID(),
				}
			}
		}
	}

	return accountKeyDetails{
		keyVaultKeyId: "",
		managedHsmId:  "",
	}
}
