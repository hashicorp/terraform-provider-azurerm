// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/packetcorecontrolplane"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/packetcoredataplane"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PacketCoreDataPlaneModel struct {
	Name                                  string            `tfschema:"name"`
	MobileNetworkPacketCoreControlPlaneId string            `tfschema:"mobile_network_packet_core_control_plane_id"`
	Location                              string            `tfschema:"location"`
	Tags                                  map[string]string `tfschema:"tags"`
	UserPlaneAccessIPv4Address            string            `tfschema:"user_plane_access_ipv4_address"`
	UserPlaneAccessIPv4Gateway            string            `tfschema:"user_plane_access_ipv4_gateway"`
	UserPlaneAccessIPv4Subnet             string            `tfschema:"user_plane_access_ipv4_subnet"`
	UserPlaneAccessName                   string            `tfschema:"user_plane_access_name"`
}

type PacketCoreDataPlaneResource struct{}

var _ sdk.ResourceWithUpdate = PacketCoreDataPlaneResource{}

func (r PacketCoreDataPlaneResource) ResourceType() string {
	return "azurerm_mobile_network_packet_core_data_plane"
}

func (r PacketCoreDataPlaneResource) ModelObject() interface{} {
	return &PacketCoreDataPlaneModel{}
}

func (r PacketCoreDataPlaneResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return packetcoredataplane.ValidatePacketCoreDataPlaneID
}

func (r PacketCoreDataPlaneResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"mobile_network_packet_core_control_plane_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: packetcorecontrolplane.ValidatePacketCoreControlPlaneID,
		},

		"location": commonschema.Location(),

		"user_plane_access_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"user_plane_access_ipv4_address": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsIPv4Address,
		},

		"user_plane_access_ipv4_subnet": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.CIDR,
		},

		"user_plane_access_ipv4_gateway": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsIPv4Address,
		},

		"tags": commonschema.Tags(),
	}
}

func (r PacketCoreDataPlaneResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r PacketCoreDataPlaneResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PacketCoreDataPlaneModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.PacketCoreDataPlaneClient
			packetCoreControlPlaneId, err := packetcorecontrolplane.ParsePacketCoreControlPlaneID(model.MobileNetworkPacketCoreControlPlaneId)
			if err != nil {
				return err
			}

			id := packetcoredataplane.NewPacketCoreDataPlaneID(packetCoreControlPlaneId.SubscriptionId, packetCoreControlPlaneId.ResourceGroupName, packetCoreControlPlaneId.PacketCoreControlPlaneName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			dataPlane := packetcoredataplane.PacketCoreDataPlane{
				Location: location.Normalize(model.Location),
				Properties: packetcoredataplane.PacketCoreDataPlanePropertiesFormat{
					UserPlaneAccessInterface: packetcoredataplane.InterfaceProperties{},
				},
				Tags: &model.Tags,
			}

			if model.UserPlaneAccessName != "" {
				dataPlane.Properties.UserPlaneAccessInterface.Name = &model.UserPlaneAccessName
			}

			if model.UserPlaneAccessIPv4Address != "" {
				dataPlane.Properties.UserPlaneAccessInterface.IPv4Address = &model.UserPlaneAccessIPv4Address
			}

			if model.UserPlaneAccessIPv4Subnet != "" {
				dataPlane.Properties.UserPlaneAccessInterface.IPv4Subnet = &model.UserPlaneAccessIPv4Subnet
			}

			if model.UserPlaneAccessIPv4Gateway != "" {
				dataPlane.Properties.UserPlaneAccessInterface.IPv4Gateway = &model.UserPlaneAccessIPv4Gateway
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, dataPlane); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r PacketCoreDataPlaneResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.PacketCoreDataPlaneClient

			id, err := packetcoredataplane.ParsePacketCoreDataPlaneID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var plan PacketCoreDataPlaneModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			model := *resp.Model

			if metadata.ResourceData.HasChange("user_plane_access_name") {
				model.Properties.UserPlaneAccessInterface.Name = &plan.UserPlaneAccessName
			}

			if metadata.ResourceData.HasChange("user_plane_access_ipv4_address") {
				model.Properties.UserPlaneAccessInterface.IPv4Address = &plan.UserPlaneAccessIPv4Address
			}

			if metadata.ResourceData.HasChange("user_plane_access_ipv4_subnet") {
				model.Properties.UserPlaneAccessInterface.IPv4Subnet = &plan.UserPlaneAccessIPv4Subnet
			}

			if metadata.ResourceData.HasChange("user_plane_access_ipv4_gateway") {
				model.Properties.UserPlaneAccessInterface.IPv4Gateway = &plan.UserPlaneAccessIPv4Gateway
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = &plan.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r PacketCoreDataPlaneResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.PacketCoreDataPlaneClient

			id, err := packetcoredataplane.ParsePacketCoreDataPlaneID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := PacketCoreDataPlaneModel{
				Name:                                  id.PacketCoreDataPlaneName,
				MobileNetworkPacketCoreControlPlaneId: packetcorecontrolplane.NewPacketCoreControlPlaneID(id.SubscriptionId, id.ResourceGroupName, id.PacketCoreControlPlaneName).ID(),
			}

			if resp.Model != nil {
				state.Location = location.Normalize(resp.Model.Location)

				props := resp.Model.Properties

				state.UserPlaneAccessIPv4Address = pointer.From(props.UserPlaneAccessInterface.IPv4Address)
				state.UserPlaneAccessIPv4Gateway = pointer.From(props.UserPlaneAccessInterface.IPv4Gateway)
				state.UserPlaneAccessIPv4Subnet = pointer.From(props.UserPlaneAccessInterface.IPv4Subnet)
				state.UserPlaneAccessName = pointer.From(props.UserPlaneAccessInterface.Name)

				if resp.Model.Tags != nil {
					state.Tags = *resp.Model.Tags
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r PacketCoreDataPlaneResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.PacketCoreDataPlaneClient

			id, err := packetcoredataplane.ParsePacketCoreDataPlaneID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// a workaround for that some child resources may still exist for seconds before it fully deleted.
			// tracked on https://github.com/Azure/azure-rest-api-specs/issues/22691
			// it will cause the error "Can not delete resource before nested resources are deleted."
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("could not retrieve context deadline for %s", id.ID())
			}
			stateConf := &pluginsdk.StateChangeConf{
				Delay:   5 * time.Minute,
				Pending: []string{"409"},
				Target:  []string{"200", "202"},
				Refresh: func() (result interface{}, state string, err error) {
					resp, err := client.Delete(ctx, *id)
					if err != nil {
						if resp.HttpResponse.StatusCode == http.StatusConflict {
							return nil, "409", nil
						}
						return nil, "", err
					}
					return resp, "200", nil
				},
				MinTimeout: 15 * time.Second,
				Timeout:    time.Until(deadline),
			}

			if future, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for deleting of %s: %+v", id, err)
			} else {
				poller := future.(packetcoredataplane.DeleteOperationResponse).Poller
				if err := poller.PollUntilDone(ctx); err != nil {
					return fmt.Errorf("deleting %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}
