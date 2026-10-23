// Copyright IBM Corp. 2014, 2025
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: keyvault.ValidateNestedItemName,
		},

		"key_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KeyVaultId{}),

		"value": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ExactlyOneOf: []string{"value", "value_wo"},
		},

		"value_wo": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			WriteOnly:    true,
			RequiredWith: []string{"value_wo_version"},
			ExactlyOneOf: []string{"value", "value_wo"},
		},

		"value_wo_version": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			RequiredWith: []string{"value_wo"},
			ValidateFunc: validation.IntAtLeast(1),
		},

		"content_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"not_before_date": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"expiration_date": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"versionless_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"resource_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"resource_versionless_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsWithMaximumElements(15),
	}
}
