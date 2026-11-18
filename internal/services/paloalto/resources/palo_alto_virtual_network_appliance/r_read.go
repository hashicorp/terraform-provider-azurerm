// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_virtual_network_appliance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkvirtualappliances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkVirtualAppliances

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
