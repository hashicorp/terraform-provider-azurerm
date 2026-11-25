// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_named_value

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiManagementNamedValueSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementChildName("The name of the API Management Named Value."),

		"resource_group_name": commonschema.ResourceGroupName(),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"display_name": {
			Description:  "The display name of this API Management Named Value.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ApiManagementNamedValueDisplayName,
		},

		"value_from_key_vault": {
			Description:  "If specified, `secret` must also be set to `true`.",
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: []string{"value", "value_from_key_vault"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"secret_id": {
						Description:  "The resource ID of the Key Vault Secret.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeSecret),
					},
					"identity_client_id": {
						Description:  "The client ID of User Assigned Identity, for the API Management Service, which will be used to access the Key Vault Secret. The System Assigned Identity will be used in absence.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
			RequiredWith: []string{"secret"},
		},

		"value": {
			Description:  "The value of this API Management Named Value.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"value", "value_from_key_vault"},
		},

		"secret": {
			Description: "Whether the API Management Named Value is secret.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"tags": {
			Description: "A list of tags to be applied to the API Management Named Value.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}
