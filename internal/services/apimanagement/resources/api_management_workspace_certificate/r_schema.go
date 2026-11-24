// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_certificate

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the API Management Workspace Certificate.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,78}[a-zA-Z0-9])?$`),
				"The `name` must be 1–80 characters, using only letters, numbers, or hyphens, and not starting or ending with a hyphen.",
			),
		},

		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&workspace.WorkspaceId{}, "API Management Workspace"),

		"certificate_data_base64": {
			Description:  "The base64-encoded string containing the certificate in PKCS#12 (.pfx) format.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ExactlyOneOf: []string{"certificate_data_base64", "key_vault_secret_id"},
			ValidateFunc: validation.StringIsBase64,
		},

		"key_vault_secret_id": {
			Description:  "The ID of the Key Vault Secret.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeSecret),
			ExactlyOneOf: []string{"certificate_data_base64", "key_vault_secret_id"},
		},

		"password": {
			Description:  "The password used to access the `certificate_data_base64`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			RequiredWith: []string{"certificate_data_base64"},
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"user_assigned_identity_client_id": {
			Description:  "The client ID of user-assigned identity to be used for accessing the `key_vault_secret_id`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
			RequiredWith: []string{"key_vault_secret_id"},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"expiration": {
			Description: "The expiration date of the API Management Workspace Certificate.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"subject": {
			Description: "The subject name of the API Management Workspace Certificate.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"thumbprint": {
			Description: "The thumbprint of the API Management Workspace Certificate.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
