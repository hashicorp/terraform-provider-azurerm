// Copyright IBM Corp. 2014, 2025
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
		"name": schemaz.SchemaApiManagementChildName(),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"authorization_endpoint": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"authorization_methods": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(authorizationserver.PossibleValuesForAuthorizationMethod(), false),
			},
			Set: pluginsdk.HashString,
		},

		"client_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"client_registration_endpoint": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"grant_types": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(authorizationserver.PossibleValuesForGrantType(), false),
			},
			Set: pluginsdk.HashString,
		},

		// Optional
		"bearer_token_sending_methods": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(authorizationserver.PossibleValuesForBearerTokenSendingMethod(), false),
			},
			Set: pluginsdk.HashString,
		},

		"client_authentication_method": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(authorizationserver.PossibleValuesForClientAuthenticationMethod(), false),
			},
			Set: pluginsdk.HashString,
		},

		"client_secret": {
			Type:      pluginsdk.TypeString,
			Optional:  true,
			Sensitive: true,
		},

		"default_scope": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"resource_owner_username": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"resource_owner_password": {
			Type:      pluginsdk.TypeString,
			Optional:  true,
			Sensitive: true,
		},

		"support_state": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"token_body_parameter": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"value": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"token_endpoint": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}
