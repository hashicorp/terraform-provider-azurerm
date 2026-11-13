// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute_fleet

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-11-01/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ComputeFleet.ComputeFleetClient

			id, err := fleets.ParseFleetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ComputeFleetResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			properties := existing.Model

			// API requires `osProfile.adminPassword` when updating resource but the GET API does not return the sensitive data `osProfile.adminPassword`
			if props := properties.Properties; props != nil {
				if v := props.ComputeProfile.BaseVirtualMachineProfile.OsProfile; v != nil {
					if len(model.VirtualMachineProfile[0].OsProfile[0].LinuxConfiguration) > 0 {
						v.AdminPassword = pointer.To(model.VirtualMachineProfile[0].OsProfile[0].LinuxConfiguration[0].AdminPassword)
					}
					if len(model.VirtualMachineProfile[0].OsProfile[0].WindowsConfiguration) > 0 {
						v.AdminPassword = pointer.To(model.VirtualMachineProfile[0].OsProfile[0].WindowsConfiguration[0].AdminPassword)
					}
				}
			}

			if metadata.ResourceData.HasChange("identity") {
				identityValue, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				properties.Identity = identityValue
			}

			if metadata.ResourceData.HasChange("plan") {
				properties.Plan = r.expandPlanModel(model.Plan)
			}

			if metadata.ResourceData.HasChange("on_demand_capacity") {
				properties.Properties.RegularPriorityProfile = r.expandOnDemandCapacityModel(model.OnDemandCapacity)
			}

			if metadata.ResourceData.HasChange("spot_capacity") {
				properties.Properties.SpotPriorityProfile = r.expandSpotCapacityModel(model.SpotCapacity)
			}

			if metadata.ResourceData.HasChange("virtual_machine_sizes_profile") {
				properties.Properties.VMSizesProfile = pointer.From(r.expandVirtualMachineSizeProfileModel(model.VirtualMachineSizesProfile, metadata))
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = pointer.To(model.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
