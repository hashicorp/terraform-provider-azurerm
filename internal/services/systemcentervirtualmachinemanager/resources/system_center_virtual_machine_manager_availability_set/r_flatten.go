// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package system_center_virtual_machine_manager_availability_set

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/availabilitysets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/vmmservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *availabilitysets.AvailabilitySetId, model *availabilitysets.AvailabilitySet) error {
	state := SystemCenterVirtualMachineManagerAvailabilitySetModel{
		Name:              id.AvailabilitySetName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)
		state.CustomLocationId = pointer.From(model.ExtendedLocation.Name)
		state.Tags = pointer.From(model.Tags)

		if model.Properties != nil {
			scvmmServerId, err := vmmservers.ParseVMmServerIDInsensitively(pointer.From(model.Properties.VMmServerId))
			if err != nil {
				return err
			}
			state.SystemCenterVirtualMachineManagerServerId = scvmmServerId.ID()
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
