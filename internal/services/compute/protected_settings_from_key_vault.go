// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachinescalesets"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/compute/2023-03-01/compute"
)

func protectedSettingsFromKeyVaultSchema(conflictsWithProtectedSettings bool) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		ConflictsWith: func() []string {
			if conflictsWithProtectedSettings {
				return []string{"protected_settings"}
			}
			return []string{}
		}(),
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"secret_url": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: keyVaultValidate.NestedItemId,
				},

				"source_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),
			},
		},
	}
}

func expandProtectedSettingsFromKeyVault(input []interface{}) *virtualmachinescalesets.KeyVaultSecretReference {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &virtualmachinescalesets.KeyVaultSecretReference{
		SecretUrl: v["secret_url"].(string),
		SourceVault: virtualmachinescalesets.SubResource{
			Id: utils.String(v["source_vault_id"].(string)),
		},
	}
}

func flattenProtectedSettingsFromKeyVault(input *compute.KeyVaultSecretReference) []interface{} {
	if input == nil {
		return nil
	}

	secretUrl := ""
	if input.SecretURL != nil {
		secretUrl = *input.SecretURL
	}

	sourceVaultId := ""
	if input.SourceVault != nil && input.SourceVault.ID != nil {
		sourceVaultId = *input.SourceVault.ID
	}

	return []interface{}{
		map[string]interface{}{
			"secret_url":      secretUrl,
			"source_vault_id": sourceVaultId,
		},
	}
}
