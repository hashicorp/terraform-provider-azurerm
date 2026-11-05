// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_kubernetes_provisioned_cluster

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *connectedclusters.ConnectedClusterId, model *connectedclusters.ConnectedCluster) error {
	state := ArcKubernetesProvisionedClusterModel{
		Name:              id.ConnectedClusterName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Identity = identity.FlattenSystemAssignedToModel(&model.Identity)
		state.Location = location.Normalize(model.Location)
		state.Tags = pointer.From(model.Tags)

		props := model.Properties
		state.AzureActiveDirectory = flattenArcKubernetesClusterAadProfile(props.AadProfile)
		state.AgentVersion = pointer.From(props.AgentVersion)
		state.Distribution = pointer.From(props.Distribution)
		state.Infrastructure = pointer.From(props.Infrastructure)
		state.KubernetesVersion = pointer.From(props.KubernetesVersion)
		state.Offering = pointer.From(props.Offering)
		state.TotalCoreCount = pointer.From(props.TotalCoreCount)
		state.TotalNodeCount = pointer.From(props.TotalNodeCount)

		arcAgentAutoUpgradeEnabled := true
		arcAgentdesiredVersion := ""
		if arcAgentProfile := props.ArcAgentProfile; arcAgentProfile != nil {
			arcAgentdesiredVersion = pointer.From(arcAgentProfile.DesiredAgentVersion)
			if arcAgentProfile.AgentAutoUpgrade != nil && *arcAgentProfile.AgentAutoUpgrade == connectedclusters.AutoUpgradeOptionsDisabled {
				arcAgentAutoUpgradeEnabled = false
			}
		}
		state.ArcAgentAutoUpgradeEnabled = arcAgentAutoUpgradeEnabled
		state.ArcAgentDesiredVersion = arcAgentdesiredVersion
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func flattenArcKubernetesClusterAadProfile(input *connectedclusters.AadProfile) []AzureActiveDirectoryModel {
	if input == nil || (input.EnableAzureRBAC == nil && input.AdminGroupObjectIDs == nil && input.TenantID == nil) {
		return make([]AzureActiveDirectoryModel, 0)
	}

	return []AzureActiveDirectoryModel{
		{
			AzureRbacEnabled:    pointer.From(input.EnableAzureRBAC),
			AdminGroupObjectIds: pointer.From(input.AdminGroupObjectIDs),
			TenantId:            pointer.From(input.TenantID),
		},
	}
}
