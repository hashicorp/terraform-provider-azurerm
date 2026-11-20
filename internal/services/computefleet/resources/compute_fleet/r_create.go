// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute_fleet

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-11-01/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ComputeFleet.ComputeFleetClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ComputeFleetResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := fleets.NewFleetID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := fleets.Fleet{
				Location: location.Normalize(model.Location),
				Plan:     r.expandPlanModel(model.Plan),
				Properties: &fleets.FleetProperties{
					RegularPriorityProfile: r.expandOnDemandCapacityModel(model.OnDemandCapacity),
					SpotPriorityProfile:    r.expandSpotCapacityModel(model.SpotCapacity),
					VMSizesProfile:         pointer.From(r.expandVirtualMachineSizeProfileModel(model.VirtualMachineSizesProfile, metadata)),
				},
			}

			if model.Tags != nil {
				properties.Tags = pointer.To(model.Tags)
			}

			if model.Zones != nil {
				properties.Zones = pointer.To(model.Zones)
			}

			expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}
			properties.Identity = expandedIdentity

			computeProfile := fleets.ComputeProfile{
				AdditionalVirtualMachineCapabilities: r.expandAdditionalCapabilities(model.AdditionalCapabilities),
				PlatformFaultDomainCount:             pointer.To(model.PlatformFaultDomainCount),
			}
			if model.ComputeApiVersion != "" {
				computeProfile.ComputeApiVersion = pointer.To(model.ComputeApiVersion)
			}

			baseVirtualMachineProfileValue, err := r.expandVirtualMachineProfileModel(model.VirtualMachineProfile, &model, metadata.ResourceData)
			if err != nil {
				return err
			}
			computeProfile.BaseVirtualMachineProfile = pointer.From(baseVirtualMachineProfileValue)
			properties.Properties.ComputeProfile = computeProfile

			if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}
