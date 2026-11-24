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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PrivateDNSResolverVirtualNetworkLinkModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.PrivateDnsResolver.VirtualNetworkLinksClient
			dnsForwardingRulesetId, err := dnsforwardingrulesets.ParseDnsForwardingRulesetID(model.DnsForwardingRulesetId)
			if err != nil {
				return err
			}

			id := virtualnetworklinks.NewVirtualNetworkLinkID(dnsForwardingRulesetId.SubscriptionId, dnsForwardingRulesetId.ResourceGroupName, dnsForwardingRulesetId.DnsForwardingRulesetName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := &virtualnetworklinks.VirtualNetworkLink{
				Properties: virtualnetworklinks.VirtualNetworkLinkProperties{
					Metadata: &model.Metadata,
					VirtualNetwork: virtualnetworklinks.SubResource{
						Id: model.VirtualNetworkId,
					},
				},
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, *properties, virtualnetworklinks.CreateOrUpdateOperationOptions{}, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}
