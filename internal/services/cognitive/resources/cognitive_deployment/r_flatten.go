// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_deployment

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/deployments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *deployments.DeploymentId, model *deployments.Deployment) error {
	state := cognitiveDeploymentModel{
		Name:               id.DeploymentName,
		CognitiveAccountId: cognitiveservicesaccounts.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName).ID(),
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	if properties := model.Properties; properties != nil {
		state.Model = flattenDeploymentModelModel(properties.Model)

		state.DynamicThrottlingEnabled = pointer.From(properties.DynamicThrottlingEnabled)
		state.RaiPolicyName = pointer.From(properties.RaiPolicyName)
		state.SpilloverDeploymentName = pointer.From(properties.SpilloverDeploymentName)
		state.VersionUpgradeOption = pointer.FromEnum(properties.VersionUpgradeOption)
	}
	if sku := flattenDeploymentSkuModel(model.Sku); sku != nil {
		state.Sku = sku
	}
	return metadata.Encode(&state)
}

func flattenDeploymentModelModel(input *deployments.DeploymentModel) []DeploymentModelModel {
	outputList := make([]DeploymentModelModel, 0, 1)
	if input == nil {
		return outputList
	}

	output := DeploymentModelModel{}
	output.Format = pointer.From(input.Format)
	output.Name = pointer.From(input.Name)
	output.Version = pointer.From(input.Version)

	return append(outputList, output)
}

func flattenDeploymentSkuModel(input *deployments.Sku) []DeploymentSkuModel {
	if input == nil {
		return []DeploymentSkuModel{}
	}
	output := DeploymentSkuModel{
		Name: input.Name,
	}
	if input.Capacity != nil {
		output.Capacity = *input.Capacity
	}
	if input.Tier != nil {
		output.Tier = string(*input.Tier)
	}
	if input.Size != nil {
		output.Size = *input.Size
	}
	if input.Family != nil {
		output.Family = *input.Family
	}
	return []DeploymentSkuModel{output}
}
