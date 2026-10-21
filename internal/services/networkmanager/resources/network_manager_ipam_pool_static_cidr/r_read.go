// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_ipam_pool_static_cidr

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/staticcidrs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.StaticCidrs

			id, err := staticcidrs.ParseStaticCidrID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			ipamPoolId := staticcidrs.NewIPamPoolID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.IpamPoolName).ID()
			schema := ManagerIpamPoolStaticCidrResourceModel{
				Name:       id.StaticCidrName,
				IpamPoolId: ipamPoolId,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					schema.AddressPrefixes = pointer.From(props.AddressPrefixes)
					schema.NumberOfIPAddressesToAllocate = pointer.From(props.NumberOfIPAddressesToAllocate)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}
