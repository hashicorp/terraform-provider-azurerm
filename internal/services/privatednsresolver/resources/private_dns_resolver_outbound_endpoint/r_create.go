// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_resolver_outbound_endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsresolvers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/outboundendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PrivateDNSResolverOutboundEndpointModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.PrivateDnsResolver.OutboundEndpointsClient
			dnsResolverId, err := dnsresolvers.ParseDnsResolverID(model.PrivateDNSResolverId)
			if err != nil {
				return err
			}

			id := outboundendpoints.NewOutboundEndpointID(dnsResolverId.SubscriptionId, dnsResolverId.ResourceGroupName, dnsResolverId.DnsResolverName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := &outboundendpoints.OutboundEndpoint{
				Location: location.Normalize(model.Location),
				Properties: outboundendpoints.OutboundEndpointProperties{
					Subnet: outboundendpoints.SubResource{
						Id: model.SubnetId,
					},
				},
				Tags: &model.Tags,
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, *properties, outboundendpoints.CreateOrUpdateOperationOptions{}, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}
