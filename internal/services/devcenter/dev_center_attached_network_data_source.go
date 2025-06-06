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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = DevCenterAttachedNetworkDataSource{}

type DevCenterAttachedNetworkDataSource struct{}

type DevCenterAttachedNetworkDataSourceModel struct {
	Name                string `tfschema:"name"`
	DevCenterId         string `tfschema:"dev_center_id"`
	NetworkConnectionId string `tfschema:"network_connection_id"`
}

func (DevCenterAttachedNetworkDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"dev_center_id": commonschema.ResourceIDReferenceRequired(&devcenters.DevCenterId{}),
	}
}

func (DevCenterAttachedNetworkDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"network_connection_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (DevCenterAttachedNetworkDataSource) ModelObject() interface{} {
	return &DevCenterAttachedNetworkDataSourceModel{}
}

func (DevCenterAttachedNetworkDataSource) ResourceType() string {
	return "azurerm_dev_center_attached_network"
}

func (r DevCenterAttachedNetworkDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.AttachedNetworkConnections
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DevCenterAttachedNetworkDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			devCenterId, err := devcenters.ParseDevCenterID(state.DevCenterId)
			if err != nil {
				return err
			}

			id := attachednetworkconnections.NewDevCenterAttachedNetworkID(subscriptionId, devCenterId.ResourceGroupName, devCenterId.DevCenterName, state.Name)

			resp, err := client.AttachedNetworksGetByDevCenter(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			state.Name = id.AttachedNetworkName
			state.DevCenterId = attachednetworkconnections.NewDevCenterID(id.SubscriptionId, id.ResourceGroupName, id.DevCenterName).ID()

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.NetworkConnectionId = props.NetworkConnectionId
				}
			}

			return metadata.Encode(&state)
		},
	}
}
