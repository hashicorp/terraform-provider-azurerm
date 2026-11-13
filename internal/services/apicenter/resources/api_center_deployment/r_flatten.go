// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_deployment

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/deployments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *deployments.DeploymentId, model *deployments.Deployment) error {
	state := ApiCenterDeploymentModel{
		Name: id.DeploymentName,
		// Id: deployments.NewDeploymentID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceName, id.ApiName, id.DeploymentName).ID(),
	}

	serviceId := services.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName)

	if props := model.Properties; props != nil {
		customProperties, err := helpers.FlattenCustomProperties(props.CustomProperties)
		if err != nil {
			return err
		}
		state.CustomProperties = customProperties

		scopedDefinitionId := fmt.Sprintf("%s%s", serviceId.ID(), pointer.From(props.DefinitionId))
		scopedEnvironmentId := fmt.Sprintf("%s%s", serviceId.ID(), pointer.From(props.EnvironmentId))

		state.DefinitionId = scopedDefinitionId
		state.Description = pointer.From(props.Description)
		state.EnvironmentId = scopedEnvironmentId
		state.ServerRuntimeUri = flattenDeploymentServer(props.Server)
		state.State = metadata.ResourceData.Get("state").(string)
		state.Title = pointer.From(props.Title)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func flattenDeploymentServer(input *deployments.DeploymentServer) []string {
	if input == nil {
		return []string{}
	}
	return pointer.From(input.RuntimeUri)
}
