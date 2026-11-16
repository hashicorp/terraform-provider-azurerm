// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_resolver_virtual_network_link

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsforwardingrulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/virtualnetworklinks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.VirtualNetworkLinksClient

			id, err := virtualnetworklinks.ParseVirtualNetworkLinkID(metadata.ResourceData.Id())
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

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := PrivateDNSResolverVirtualNetworkLinkModel{
				Name:                   id.VirtualNetworkLinkName,
				DnsForwardingRulesetId: dnsforwardingrulesets.NewDnsForwardingRulesetID(id.SubscriptionId, id.ResourceGroupName, id.DnsForwardingRulesetName).ID(),
			}

			properties := &model.Properties
			if properties.Metadata != nil {
				state.Metadata = *properties.Metadata
			}

			state.VirtualNetworkId = properties.VirtualNetwork.Id

			return metadata.Encode(&state)
		},
	}
}
