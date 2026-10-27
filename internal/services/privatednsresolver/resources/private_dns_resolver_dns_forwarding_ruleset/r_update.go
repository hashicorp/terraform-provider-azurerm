// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_resolver_dns_forwarding_ruleset

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsforwardingrulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.DnsForwardingRulesetsClient

			id, err := dnsforwardingrulesets.ParseDnsForwardingRulesetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model PrivateDNSResolverDnsForwardingRulesetModel
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

			if metadata.ResourceData.HasChange("private_dns_resolver_outbound_endpoint_ids") {
				dnsResolverOutboundEndpointsValue := expandDnsResolverOutboundEndpoints(model.DnsResolverOutboundEndpoints)

				if dnsResolverOutboundEndpointsValue != nil {
					properties.Properties.DnsResolverOutboundEndpoints = *dnsResolverOutboundEndpointsValue
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties, dnsforwardingrulesets.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
