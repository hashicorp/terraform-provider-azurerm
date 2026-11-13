// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_deployment

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apidefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/deployments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Description:  "The name of the deployment.",
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.Name(),
		},

		"api_definition_id": {
			Type:         pluginsdk.TypeString,
			Description:  "The ID of the API definition to deploy.",
			Required:     true,
			ValidateFunc: apidefinitions.ValidateDefinitionID,
		},

		"environment_id": {
			Type:         pluginsdk.TypeString,
			Description:  "The ID of the environment for this deployment.",
			Required:     true,
			ValidateFunc: environments.ValidateEnvironmentID,
		},

		"server_runtime_uri": {
			Type:        pluginsdk.TypeSet,
			Description: "A set of server runtime URLs for this deployment.",
			Required:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				Description:  "A URL for a server runtime.",
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
		},

		"custom_properties": helpers.CustomPropertiesSchema(),

		"description": {
			Type:         pluginsdk.TypeString,
			Description:  "A description of the deployment.",
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 500),
		},

		"display_name": {
			Type:        pluginsdk.TypeString,
			Description: "The display name of the deployment.",
			Optional:    true,
			// Note: O+C if display name is not set, it will default to the name of the deployment
			Computed:     true,
			ValidateFunc: validation.StringLenBetween(1, 50),
		},

		"state": {
			Type:         pluginsdk.TypeString,
			Description:  "The state of the deployment.",
			Optional:     true,
			ValidateFunc: validation.StringInSlice(deployments.PossibleValuesForDeploymentState(), false),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
