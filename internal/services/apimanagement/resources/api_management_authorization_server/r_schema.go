// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_authorization_server

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/authorizationserver"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiManagementAuthorizationServerSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementChildName("The name of this Authorization Server."),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"authorization_endpoint": {
			Description:  "The OAuth Authorization Endpoint.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"authorization_methods": {
			Description: "The HTTP Verbs supported by the Authorization Endpoint.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(authorizationserver.PossibleValuesForAuthorizationMethod(), false),
			},
			Set: pluginsdk.HashString,
		},

		"client_id": {
			Description:  "The Client/App ID registered with this Authorization Server.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"client_registration_endpoint": {
			Description:  "The URI of page where Client/App Registration is performed for this Authorization Server.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"display_name": {
			Description:  "The user-friendly name of this Authorization Server.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"grant_types": {
			Description: "Form of Authorization Grants required when requesting an Access Token.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(authorizationserver.PossibleValuesForGrantType(), false),
			},
			Set: pluginsdk.HashString,
		},

		// Optional
		"bearer_token_sending_methods": {
			Description: "The mechanism by which Access Tokens are passed to the API.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(authorizationserver.PossibleValuesForBearerTokenSendingMethod(), false),
			},
			Set: pluginsdk.HashString,
		},

		"client_authentication_method": {
			Description: "The Authentication Methods supported by the Token endpoint of this Authorization Server.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(authorizationserver.PossibleValuesForClientAuthenticationMethod(), false),
			},
			Set: pluginsdk.HashString,
		},

		"client_secret": {
			Description: "The Client/App Secret registered with this Authorization Server.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Sensitive:   true,
		},

		"default_scope": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"description": {
			Description: "A description of the Authorization Server, which may contain HTML formatting tags.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"resource_owner_username": {
			Description: "The username associated with the Resource Owner.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"resource_owner_password": {
			Description: "The password associated with the Resource Owner.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Sensitive:   true,
		},

		"support_state": {
			Description: "Does this Authorization Server support State? If this is set to `true` the client may use the state parameter to raise protocol security.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"token_body_parameter": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The Name of the Parameter.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"value": {
						Description:  "The Value of the Parameter.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"token_endpoint": {
			Description: "The OAuth Token Endpoint.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},
	}
}
