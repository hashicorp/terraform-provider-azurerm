// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_update_run

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/updateruns"
)

func flattenKubernetesFleetUpdateRunManagedClusterUpdate(input updateruns.ManagedClusterUpdate) []KubernetesFleetUpdateRunResourceManagedClusterUpdateSchema {
	return []KubernetesFleetUpdateRunResourceManagedClusterUpdateSchema{
		{
			NodeImageSelection: flattenKubernetesFleetUpdateRunNodeImageSelection(input.NodeImageSelection),
			Upgrade:            flattenKubernetesFleetUpdateRunUpgrade(input.Upgrade),
		},
	}
}

func flattenKubernetesFleetUpdateRunNodeImageSelection(input *updateruns.NodeImageSelection) []KubernetesFleetUpdateRunResourceManagedClusterUpdateNodeImageSelectionSchema {
	if input == nil {
		return []KubernetesFleetUpdateRunResourceManagedClusterUpdateNodeImageSelectionSchema{}
	}
	return []KubernetesFleetUpdateRunResourceManagedClusterUpdateNodeImageSelectionSchema{
		{
			Type: string(input.Type),
		},
	}
}

func flattenKubernetesFleetUpdateRunUpgrade(input updateruns.ManagedClusterUpgradeSpec) []KubernetesFleetUpdateRunResourceManagedClusterUpdateUpgradeSchema {
	return []KubernetesFleetUpdateRunResourceManagedClusterUpdateUpgradeSchema{
		{
			Type:              string(input.Type),
			KubernetesVersion: pointer.From(input.KubernetesVersion),
		},
	}
}

func flattenKubernetesFleetUpdateRunStage(input []updateruns.UpdateStage) []KubernetesFleetUpdateRunResourceUpdateStageSchema {
	output := make([]KubernetesFleetUpdateRunResourceUpdateStageSchema, 0)
	for _, stage := range input {
		output = append(output, KubernetesFleetUpdateRunResourceUpdateStageSchema{
			Name:                    stage.Name,
			AfterStageWaitInSeconds: pointer.From(stage.AfterStageWaitInSeconds),
			Group:                   flattenKubernetesFleetUpdateRunGroup(stage.Groups),
		})
	}
	return output
}

func flattenKubernetesFleetUpdateRunGroup(input *[]updateruns.UpdateGroup) []KubernetesFleetUpdateRunResourceUpdateGroupSchema {
	output := make([]KubernetesFleetUpdateRunResourceUpdateGroupSchema, 0)
	if input == nil {
		return output
	}

	for _, group := range *input {
		output = append(output, KubernetesFleetUpdateRunResourceUpdateGroupSchema{
			Name: group.Name,
		})
	}

	return output
}
