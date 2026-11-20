// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_custom_domain

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managedenvironments.ValidateManagedEnvironmentID,
			Description:  "The Container App Managed Environment ID to configure this Custom Domain on.",
		},

		"certificate_blob_base64": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsBase64,
			ExactlyOneOf: []string{"certificate_key_vault", "certificate_blob_base64"},
			RequiredWith: []string{"certificate_password"},
			Description:  "The Custom Domain Certificate Private Key as a base64 encoded PFX or PEM.",
		},

		"certificate_password": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			RequiredWith: []string{"certificate_blob_base64"},
			Description:  "The Custom Domain Certificate password.",
		},

		"certificate_key_vault": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: []string{"certificate_key_vault", "certificate_blob_base64"},
			Description:  "A `certificate_key_vault` block as defined below.",
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"identity": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "System",
						ValidateFunc: validation.Any(
							validation.StringInSlice([]string{"System"}, false),
							commonids.ValidateUserAssignedIdentityID,
						),
						Description: "The identity to use for accessing the Key Vault. Use `System` for system-assigned identity or a User Assigned Identity ID.",
					},
					"key_vault_secret_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersionless, keyvault.NestedItemTypeSecret),
						Description:  "The URL of the Key Vault secret containing the certificate.",
					},
				},
			},
		},

		"dns_suffix": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The Custom Domain DNS suffix for this Container App Environment.",
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
