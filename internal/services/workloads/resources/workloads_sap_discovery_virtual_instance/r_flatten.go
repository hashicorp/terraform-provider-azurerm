// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package workloads_sap_discovery_virtual_instance

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2024-09-01/sapvirtualinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (Resource) flatten(metadata sdk.ResourceMetaData, id *sapvirtualinstances.SapVirtualInstanceId, model *sapvirtualinstances.SAPVirtualInstance) error {
	state := WorkloadsSAPDiscoveryVirtualInstanceModel{
		Name:              id.SapVirtualInstanceName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)

		identity, err := identity.FlattenUserAssignedMapToModel(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		state.Identity = pointer.From(identity)

		if props := model.Properties; props != nil {
			state.Environment = string(props.Environment)
			state.ManagedResourcesNetworkAccessType = pointer.FromEnum(props.ManagedResourcesNetworkAccessType)
			state.SapProduct = string(props.SapProduct)
			state.Tags = pointer.From(model.Tags)

			if config := props.Configuration; config != nil {
				if v, ok := config.(sapvirtualinstances.DiscoveryConfiguration); ok {
					state.CentralServerVmId = pointer.From(v.CentralServerVMId)
					state.ManagedStorageAccountName = pointer.From(v.ManagedRgStorageAccountName)
				}
			}

			if v := props.ManagedResourceGroupConfiguration; v != nil {
				state.ManagedResourceGroupName = pointer.From(v.Name)
			}
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
