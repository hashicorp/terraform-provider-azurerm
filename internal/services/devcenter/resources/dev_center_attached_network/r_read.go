// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_attached_network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/attachednetworkconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
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
