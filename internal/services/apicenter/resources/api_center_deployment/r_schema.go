// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_deployment

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apidefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the deployment.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.Name(),
		},

		"api_center_api_definition_id": {
			Description:  "The ID of the API definition to deploy.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: apidefinitions.ValidateDefinitionID,
		},

		"api_center_environment_id": {
			Description:  "The ID of the environment for this deployment.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: environments.ValidateEnvironmentID,
		},

		"server_runtime_url": {
			Description: "A set of server runtime URLs for this deployment.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			Elem: &pluginsdk.Schema{
				Description:  "A URL for a server runtime.",
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
		},

		"custom_properties_json": helpers.CustomPropertiesSchema(),

		"description": {
			Description:  "A description of the deployment.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 500),
		},

		"display_name": {
			Description: "The display name of the deployment.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C if display name is not set, it will default to the name of the deployment
			Computed:     true,
			ValidateFunc: validation.StringLenBetween(1, 50),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
