// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_identity_provider_aadb2c

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiManagementIdentityProviderAadb2cSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"client_id": {
			Description:  "Client ID of the Application in your B2C tenant.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},

		"client_secret": {
			Description:  "Client secret of the Application in your B2C tenant.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		// For AADB2C identity providers, `allowed_tenants` must specify exactly one tenant
		"allowed_tenant": {
			Description:  "The allowed AAD tenant, usually your B2C tenant domain.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"signin_tenant": {
			Description: "The tenant to use instead of Common when logging into Active Directory, usually your B2C tenant domain.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			// B2C tenant domains can be customized, and GUIDs might work here too
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"authority": {
			Description: "OpenID Connect discovery endpoint hostname, usually your b2clogin.com domain.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			// B2C login domains can be customized and don't necessarily end in b2clogin.com
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"signup_policy": {
			Description:  "Signup Policy Name.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"signin_policy": {
			Description:  "Signin Policy Name.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"client_library": {
			Description:  "The client library to be used in the Azure AD B2C Identity Provider.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(0, 16),
		},

		"profile_editing_policy": {
			Description:  "Profile editing Policy Name.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"password_reset_policy": {
			Description:  "Password reset Policy Name.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}
