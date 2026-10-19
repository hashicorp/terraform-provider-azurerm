// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachineextensions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesetextensions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func ProtectedSettingsFromKeyVaultSchema(conflictsWithProtectedSettings bool) *pluginsdk.Schema {
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
					ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeSecret),
				},

				"source_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),
			},
		},
	}
}

func ExpandProtectedSettingsFromKeyVault(input []any) *virtualmachineextensions.KeyVaultSecretReference {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	return &virtualmachineextensions.KeyVaultSecretReference{
		SecretURL: v["secret_url"].(string),
		SourceVault: virtualmachineextensions.SubResource{
			Id: pointer.To(v["source_vault_id"].(string)),
		},
	}
}

func ExpandProtectedSettingsFromKeyVaultVMSS(input []any) *virtualmachinescalesets.KeyVaultSecretReference {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	return &virtualmachinescalesets.KeyVaultSecretReference{
		SecretURL: v["secret_url"].(string),
		SourceVault: virtualmachinescalesets.SubResource{
			Id: pointer.To(v["source_vault_id"].(string)),
		},
	}
}

func ExpandProtectedSettingsFromKeyVaultOldVMSSExtension(input []any) *virtualmachinescalesetextensions.KeyVaultSecretReference {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	return &virtualmachinescalesetextensions.KeyVaultSecretReference{
		SecretURL: v["secret_url"].(string),
		SourceVault: virtualmachinescalesetextensions.SubResource{
			Id: pointer.To(v["source_vault_id"].(string)),
		},
	}
}

func FlattenProtectedSettingsFromKeyVault(input *virtualmachineextensions.KeyVaultSecretReference) []any {
	if input == nil {
		return []any{}
	}

	sourceVaultId := pointer.From(input.SourceVault.Id)

	return []any{
		map[string]any{
			"secret_url":      input.SecretURL,
			"source_vault_id": sourceVaultId,
		},
	}
}

func FlattenProtectedSettingsFromKeyVaultVMSS(input *virtualmachinescalesets.KeyVaultSecretReference) []any {
	if input == nil {
		return []any{}
	}

	sourceVaultId := pointer.From(input.SourceVault.Id)

	return []any{
		map[string]any{
			"secret_url":      input.SecretURL,
			"source_vault_id": sourceVaultId,
		},
	}
}

func FlattenProtectedSettingsFromKeyVaultOldVMSSExtension(input *virtualmachinescalesetextensions.KeyVaultSecretReference) []any {
	if input == nil {
		return []any{}
	}

	sourceVaultId := pointer.From(input.SourceVault.Id)

	return []any{
		map[string]any{
			"secret_url":      input.SecretURL,
			"source_vault_id": sourceVaultId,
		},
	}
}
