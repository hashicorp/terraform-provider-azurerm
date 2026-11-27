// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_mover_agent

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2023-03-01/agents"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2023-03-01/storagemovers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *agents.AgentId, model *agents.Agent) error {
	state := StorageMoverAgentResourceModel{
		Name:           id.AgentName,
		StorageMoverId: storagemovers.NewStorageMoverID(id.SubscriptionId, id.ResourceGroupName, id.StorageMoverName).ID(),
	}

	if model != nil {
		state.ArcVmUuid = model.Properties.ArcVMUuid
		state.ArcVirtualMachineId = model.Properties.ArcResourceId
		state.Description = pointer.From(model.Properties.Description)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
