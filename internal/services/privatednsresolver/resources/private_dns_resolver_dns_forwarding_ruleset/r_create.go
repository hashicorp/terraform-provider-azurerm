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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PrivateDNSResolverDnsForwardingRulesetModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.PrivateDnsResolver.DnsForwardingRulesetsClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := dnsforwardingrulesets.NewDnsForwardingRulesetID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := &dnsforwardingrulesets.DnsForwardingRuleset{
				Location:   location.Normalize(model.Location),
				Properties: dnsforwardingrulesets.DnsForwardingRulesetProperties{},
				Tags:       &model.Tags,
			}

			dnsResolverOutboundEndpointsValue := expandDnsResolverOutboundEndpoints(model.DnsResolverOutboundEndpoints)

			if dnsResolverOutboundEndpointsValue != nil {
				properties.Properties.DnsResolverOutboundEndpoints = *dnsResolverOutboundEndpointsValue
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, *properties, dnsforwardingrulesets.CreateOrUpdateOperationOptions{}, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}

func expandDnsResolverOutboundEndpoints(inputList []string) *[]dnsforwardingrulesets.SubResource {
	outputList := make([]dnsforwardingrulesets.SubResource, 0, len(inputList))
	for _, v := range inputList {
		output := dnsforwardingrulesets.SubResource{
			Id: v,
		}
		outputList = append(outputList, output)
	}

	return &outputList
}
