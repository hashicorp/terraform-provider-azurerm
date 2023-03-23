package appconfiguration

import (
	"context"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2022-05-01/configurationstores"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/sdk/1.0/appconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
func appConfigurationGetKeyRefreshFunc(ctx context.Context, client *appconfiguration.BaseClient, key, label string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.GetKeyValue(ctx, key, label, "", "", "", []string{})
		if err != nil {
			if v, ok := err.(autorest.DetailedError); ok {
				if utils.ResponseWasForbidden(autorest.Response{Response: v.Response}) {
					return "Forbidden", "Forbidden", nil
				}
			}
			return res, "Error", nil
		}

		return res, "Exists", nil
	}
}
