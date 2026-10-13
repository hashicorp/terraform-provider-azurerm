// Copyright IBM Corp. 2014, 2025
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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,78}[a-zA-Z0-9])?$`),
				"The `name` must be 1–80 characters, using only letters, numbers, or hyphens, and not starting or ending with a hyphen.",
			),
		},

		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&workspace.WorkspaceId{}),

		"certificate_data_base64": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ExactlyOneOf: []string{"certificate_data_base64", "key_vault_secret_id"},
			ValidateFunc: validation.StringIsBase64,
		},

		"key_vault_secret_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeSecret),
			ExactlyOneOf: []string{"certificate_data_base64", "key_vault_secret_id"},
		},

		"password": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			RequiredWith: []string{"certificate_data_base64"},
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"user_assigned_identity_client_id": {
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
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"subject": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"thumbprint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
