// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/attachednetworkconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devcenters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/networkconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = DevCenterAttachedNetworkResource{}

type DevCenterAttachedNetworkResource struct{}

func (r DevCenterAttachedNetworkResource) ModelObject() interface{} {
	return &DevCenterAttachedNetworkResourceModel{}
}

type DevCenterAttachedNetworkResourceModel struct {
	Name                string `tfschema:"name"`
	DevCenterId         string `tfschema:"dev_center_id"`
	NetworkConnectionId string `tfschema:"network_connection_id"`
}

func (r DevCenterAttachedNetworkResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return attachednetworkconnections.ValidateDevCenterAttachedNetworkID
}

func (r DevCenterAttachedNetworkResource) ResourceType() string {
	return "azurerm_dev_center_attached_network"
}

func (r DevCenterAttachedNetworkResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"dev_center_id": commonschema.ResourceIDReferenceRequiredForceNew(&devcenters.DevCenterId{}),

		"network_connection_id": commonschema.ResourceIDReferenceRequiredForceNew(&networkconnections.NetworkConnectionId{}),
	}
}

func (r DevCenterAttachedNetworkResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DevCenterAttachedNetworkResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.AttachedNetworkConnections
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DevCenterAttachedNetworkResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			devCenterId, err := devcenters.ParseDevCenterID(model.DevCenterId)
			if err != nil {
				return err
			}

			id := attachednetworkconnections.NewDevCenterAttachedNetworkID(subscriptionId, devCenterId.ResourceGroupName, devCenterId.DevCenterName, model.Name)

			existing, err := client.AttachedNetworksGetByDevCenter(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := attachednetworkconnections.AttachedNetworkConnection{
				Properties: &attachednetworkconnections.AttachedNetworkConnectionProperties{
					NetworkConnectionId: model.NetworkConnectionId,
				},
			}

			if err := client.AttachedNetworksCreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DevCenterAttachedNetworkResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.AttachedNetworkConnections

			id, err := attachednetworkconnections.ParseDevCenterAttachedNetworkID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.AttachedNetworksGetByDevCenter(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := DevCenterAttachedNetworkResourceModel{
				Name:        id.AttachedNetworkName,
				DevCenterId: attachednetworkconnections.NewDevCenterID(id.SubscriptionId, id.ResourceGroupName, id.DevCenterName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.NetworkConnectionId = props.NetworkConnectionId
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DevCenterAttachedNetworkResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.AttachedNetworkConnections

			id, err := attachednetworkconnections.ParseDevCenterAttachedNetworkID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.AttachedNetworksDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
