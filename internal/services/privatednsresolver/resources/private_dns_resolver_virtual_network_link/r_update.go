// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_resolver_virtual_network_link

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/virtualnetworklinks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.VirtualNetworkLinksClient

			id, err := virtualnetworklinks.ParseVirtualNetworkLinkID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model PrivateDNSResolverVirtualNetworkLinkModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("metadata") {
				properties.Properties.Metadata = &model.Metadata
			}

			if metadata.ResourceData.HasChange("virtual_network_id") {
				properties.Properties.VirtualNetwork = virtualnetworklinks.SubResource{
					Id: model.VirtualNetworkId,
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties, virtualnetworklinks.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
