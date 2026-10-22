// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_virtual_network_appliance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkvirtualappliances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkVirtualAppliances

			model := NetworkVirtualApplianceResourceModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			hubID, err := virtualwans.ParseVirtualHubID(model.VirtualHubID)
			if err != nil {
				return err
			}

			id := networkvirtualappliances.NewNetworkVirtualApplianceID(hubID.SubscriptionId, hubID.ResourceGroupName, model.Name)

			hub, err := metadata.Client.Network.VirtualWANs.VirtualHubsGet(ctx, *hubID)
			if err != nil {
				return fmt.Errorf("retrieving %s for %s: %+v", hubID, id, err)
			}
			if hub.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", hubID)
			}

			loc := location.Normalize(pointer.From(hub.Model.Location))

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id, networkvirtualappliances.DefaultGetOperationOptions())
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			props := networkvirtualappliances.NetworkVirtualAppliancePropertiesFormat{
				Delegation: &networkvirtualappliances.DelegationProperties{
					ServiceName: pointer.To("PaloAltoNetworks.Cloudngfw/firewalls"),
				},
				VirtualHub: &networkvirtualappliances.SubResource{
					Id: pointer.To(hubID.ID()),
				},
			}

			appliance := networkvirtualappliances.NetworkVirtualAppliance{
				Location:   pointer.To(loc),
				Properties: pointer.To(props),
			}

			if err = client.CreateOrUpdateCallbackThenPoll(ctx, id, appliance, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating Virtual Network Appliance for %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}
