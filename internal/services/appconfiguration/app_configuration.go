// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appconfiguration

import (
	"context"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/configurationstores"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/tombuildsstuff/kermit/sdk/appconfiguration/1.0/appconfiguration"
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
		res, err := client.GetKeyValue(ctx, key, label, "", "", "", []appconfiguration.KeyValueFields{})
		if err != nil {
			if v, ok := err.(autorest.DetailedError); ok {
				if response.WasForbidden(v.Response) {
					return "Forbidden", "Forbidden", nil
				}
				if response.WasNotFound(v.Response) {
					return "NotFound", "NotFound", nil
				}
			}
			return res, "Error", nil
		}

		return res, "Exists", nil
	}
}
