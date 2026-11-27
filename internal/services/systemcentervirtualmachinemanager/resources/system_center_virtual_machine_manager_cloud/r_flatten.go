// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package system_center_virtual_machine_manager_cloud

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/clouds"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *clouds.CloudId, model *clouds.Cloud) error {
	state := SystemCenterVirtualMachineManagerCloudModel{
		Name:              id.CloudName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)
		state.CustomLocationId = pointer.From(model.ExtendedLocation.Name)
		state.Tags = pointer.From(model.Tags)

		if model.Properties != nil {
			state.SystemCenterVirtualMachineManagerServerInventoryItemId = pointer.From(model.Properties.InventoryItemId)
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
