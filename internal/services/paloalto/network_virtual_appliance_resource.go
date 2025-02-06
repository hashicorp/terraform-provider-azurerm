// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package paloalto

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/networkvirtualappliances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/virtualwans"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetworkVirtualApplianceResource struct{}

type NetworkVirtualApplianceResourceModel struct {
	Name         string `tfschema:"name"`
	VirtualHubID string `tfschema:"virtual_hub_id"`
}

var _ sdk.Resource = NetworkVirtualApplianceResource{}

func (r NetworkVirtualApplianceResource) ResourceType() string {
	return "azurerm_palo_alto_virtual_network_appliance"
}

func (r NetworkVirtualApplianceResource) ModelObject() interface{} {
	return &NetworkVirtualApplianceResourceModel{}
}

func (r NetworkVirtualApplianceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return networkvirtualappliances.ValidateNetworkVirtualApplianceID
}

func (r NetworkVirtualApplianceResource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VirtualHubName,
		},

		"virtual_hub_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: virtualwans.ValidateVirtualHubID,
		},
	}
}

func (r NetworkVirtualApplianceResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NetworkVirtualApplianceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.Client.NetworkVirtualAppliances

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

			existing, err := client.Get(ctx, id, networkvirtualappliances.DefaultGetOperationOptions())
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
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

			if err = client.CreateOrUpdateThenPoll(ctx, id, appliance); err != nil {
				return fmt.Errorf("creating Virtual Network Appliance for %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r NetworkVirtualApplianceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.Client.NetworkVirtualAppliances

			var state NetworkVirtualApplianceResourceModel

			id, err := networkvirtualappliances.ParseNetworkVirtualApplianceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id, networkvirtualappliances.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state.Name = id.NetworkVirtualApplianceName
			if model := existing.Model; model != nil {
				if props := model.Properties; props != nil {
					if props.VirtualHub != nil {
						state.VirtualHubID = pointer.From(props.VirtualHub.Id)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r NetworkVirtualApplianceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.Client.NetworkVirtualAppliances

			id, err := networkvirtualappliances.ParseNetworkVirtualApplianceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
