// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_integration_account_certificate

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func logicAppIntegrationAccountCertificateSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Logic App Integration Account Certificate.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IntegrationAccountCertificateName(),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"integration_account_name": {
			Description:  "The name of the Logic App Integration Account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IntegrationAccountName(),
		},

		"key_vault_key": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			AtLeastOneOf: []string{"public_certificate"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_name": {
						Description:  "The name of Key Vault Key.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemName,
					},

					"key_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}, "Key Vault"),

					"key_version": {
						Description:  "The version of Key Vault Key.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"metadata": {
			Description:      "A JSON mapping of any Metadata for this Logic App Integration Account Certificate.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"public_certificate": {
			Description:  "The public certificate for the Logic App Integration Account Certificate.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			AtLeastOneOf: []string{"key_vault_key"},
		},
	}
}
