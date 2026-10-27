// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_resolver_dns_forwarding_ruleset

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsforwardingrulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.DnsForwardingRulesetsClient

			id, err := dnsforwardingrulesets.ParseDnsForwardingRulesetID(metadata.ResourceData.Id())
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

			state := PrivateDNSResolverDnsForwardingRulesetModel{
				Name:              id.DnsForwardingRulesetName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			properties := &model.Properties

			state.DnsResolverOutboundEndpoints = flattenDnsResolverOutboundEndpoints(&properties.DnsResolverOutboundEndpoints)

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenDnsResolverOutboundEndpoints(inputList *[]dnsforwardingrulesets.SubResource) []string {
	if inputList == nil {
		return []string{}
	}

	outputList := make([]string, 0, len(*inputList))
	for _, input := range *inputList {
		output := input.Id
		outputList = append(outputList, output)
	}

	return outputList
}
