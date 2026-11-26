// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_secret

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func keyVaultSecretSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Key Vault Secret.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: keyvault.ValidateNestedItemName,
		},

		"key_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KeyVaultId{}, "Key Vault"),

		"value": {
			Description:  "The value of the Key Vault Secret.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ExactlyOneOf: []string{"value", "value_wo"},
		},

		"value_wo": {
			Description:  "The value of the Key Vault Secret.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			WriteOnly:    true,
			RequiredWith: []string{"value_wo_version"},
			ExactlyOneOf: []string{"value", "value_wo"},
		},

		"value_wo_version": {
			Description:  "An integer value used to trigger an update for `value_wo`. This property should be incremented when updating `value_wo`.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
			RequiredWith: []string{"value_wo"},
		},

		"content_type": {
			Description: "The content type for the Key Vault Secret.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"not_before_date": {
			Description:  "Key not usable before the provided UTC datetime (Y-m-d'T'H:M:S'Z').",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"expiration_date": {
			Description:  "Expiration UTC datetime (Y-m-d'T'H:M:S'Z').",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"version": {
			Description: "The current version of the Key Vault Secret.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"versionless_id": {
			Description: "The Base ID of the Key Vault Secret.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"resource_id": {
			Description: "The (Versioned) ID for this Key Vault Secret. This property points to a specific version of a Key Vault Secret, as such using this won't auto-rotate values if used in other Azure Services.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"resource_versionless_id": {
			Description: "The Versionless ID of the Key Vault Secret. This property allows other Azure Services (that support it) to auto-rotate their value when the Key Vault Secret is updated.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.TagsWithMaximumElements(15),
	}
}
