package api_center_deployment

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/deployments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *deployments.DeploymentId, model *deployments.Deployment) error {
	state := ApiCenterDeploymentModel{
		Name: id.DeploymentName,
		// Id: deployments.NewDeploymentID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceName, id.ApiName, id.DeploymentName).ID(),
	}

	if props := model.Properties; props != nil {
		state.DefinitionId = pointer.From(props.DefinitionId)
		state.Description = pointer.From(props.Description)
		state.EnvironmentId = pointer.From(props.EnvironmentId)
		state.Server = flattenDeploymentServer(props.Server)
		state.State = string(pointer.From(props.State))
		state.Title = pointer.From(props.Title)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func flattenDeploymentServer(input *deployments.DeploymentServer) []DeploymentServer {
	if input == nil {
		return []DeploymentServer{}
	}
	return []DeploymentServer{
		{
			RuntimeUri: pointer.From(input.RuntimeUri),
		},
	}
}
