// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_environment

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/environments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Description:  "The name of the environment.",
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.Name(),
		},

		"service_id": {
			Type:         pluginsdk.TypeString,
			Description:  "The ID of the API Center service.",
			Required:     true,
			ForceNew:     true,
			ValidateFunc: services.ValidateServiceID,
		},

		"environment_type": {
			Type:        pluginsdk.TypeString,
			Description: "The type of the environment.",
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice(
				environments.PossibleValuesForEnvironmentKind(),
				false),
		},

		"custom_properties": helpers.CustomPropertiesSchema(),

		"description": {
			Type:         pluginsdk.TypeString,
			Description:  "A description of the environment.",
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"development_portal_uri": {
			Type:         pluginsdk.TypeString,
			Description:  "The HTTPS URI of the development portal.",
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPS,
		},

		"display_name": {
			Type:        pluginsdk.TypeString,
			Description: "The display name of the environment.",
			Optional:    true,
			// Note: O+C if display name is not set, it will default to the name of the environment
			Computed:     true,
			ValidateFunc: validation.StringLenBetween(1, 50),
		},

		"instructions": {
			Type:         pluginsdk.TypeString,
			Description:  "Instructions for using the environment.",
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"management_portal_uri": {
			Type:         pluginsdk.TypeString,
			Description:  "The HTTPS URI of the management portal.",
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPS,
		},

		"server_type": {
			Type:        pluginsdk.TypeString,
			Description: "The type of server used by the environment.",
			Optional:    true,
			ValidateFunc: validation.StringInSlice(
				environments.PossibleValuesForEnvironmentServerType(),
				false),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
