package appconfiguration

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2022-05-01/configurationstores"
)

func flattenAppConfigurationEncryption(input *configurationstores.EncryptionProperties) []interface{} {
	if input == nil || input.KeyVaultProperties == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"identity_client_id":       input.KeyVaultProperties.IdentityClientId,
			"key_vault_key_identifier": input.KeyVaultProperties.KeyIdentifier,
		},
	}
}
