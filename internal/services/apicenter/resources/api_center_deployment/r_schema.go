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
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.Name(),
		},
		"api_definition_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: apidefinitions.ValidateDefinitionID,
		},
		"custom_properties": helpers.CustomPropertiesSchema(),
		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 500),
		},
		"environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: environments.ValidateEnvironmentID,
		},
		"server_runtime_uri": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
		},
		"state": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(deployments.PossibleValuesForDeploymentState(), false),
		},
		"display_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// O+C if display name is not set, it will default to the name of the deployment
			Computed:     true,
			ValidateFunc: validation.StringLenBetween(1, 50),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
