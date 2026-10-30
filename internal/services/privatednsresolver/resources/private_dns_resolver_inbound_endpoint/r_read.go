// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_resolver_inbound_endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsresolvers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/inboundendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.InboundEndpointsClient

			id, err := inboundendpoints.ParseInboundEndpointID(metadata.ResourceData.Id())
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

			state := PrivateDNSResolverInboundEndpointModel{
				Name:                 id.InboundEndpointName,
				PrivateDNSResolverId: dnsresolvers.NewDnsResolverID(id.SubscriptionId, id.ResourceGroupName, id.DnsResolverName).ID(),
				Location:             location.Normalize(model.Location),
			}

			properties := &model.Properties

			state.IPConfigurations = flattenIPConfigurationModel(&properties.IPConfigurations)
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenIPConfigurationModel(inputList *[]inboundendpoints.IPConfiguration) []IPConfigurationModel {
	if inputList == nil {
		return []IPConfigurationModel{}
	}

	outputList := make([]IPConfigurationModel, 0, len(*inputList))
	for _, input := range *inputList {
		output := IPConfigurationModel{}

		if input.PrivateIPAddress != nil {
			output.PrivateIPAddress = *input.PrivateIPAddress
		}

		if input.PrivateIPAllocationMethod != nil {
			output.PrivateIPAllocationMethod = *input.PrivateIPAllocationMethod
		}

		output.SubnetId = input.Subnet.Id

		outputList = append(outputList, output)
	}

	return outputList
}
