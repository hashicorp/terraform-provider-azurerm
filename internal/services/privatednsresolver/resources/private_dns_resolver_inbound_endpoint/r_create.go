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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PrivateDNSResolverInboundEndpointModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.PrivateDnsResolver.InboundEndpointsClient
			dnsResolverId, err := dnsresolvers.ParseDnsResolverID(model.PrivateDNSResolverId)
			if err != nil {
				return err
			}

			id := inboundendpoints.NewInboundEndpointID(dnsResolverId.SubscriptionId, dnsResolverId.ResourceGroupName, dnsResolverId.DnsResolverName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := &inboundendpoints.InboundEndpoint{
				Location:   location.Normalize(model.Location),
				Properties: inboundendpoints.InboundEndpointProperties{},
				Tags:       &model.Tags,
			}

			iPConfigurationsValue, err := expandIPConfigurationModel(model.IPConfigurations)
			if err != nil {
				return err
			}

			if iPConfigurationsValue != nil {
				properties.Properties.IPConfigurations = *iPConfigurationsValue
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, *properties, inboundendpoints.CreateOrUpdateOperationOptions{}, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}

func expandIPConfigurationModel(inputList []IPConfigurationModel) (*[]inboundendpoints.IPConfiguration, error) {
	outputList := make([]inboundendpoints.IPConfiguration, 0, len(inputList))
	for _, v := range inputList {
		input := v
		output := inboundendpoints.IPConfiguration{}

		if input.PrivateIPAllocationMethod == inboundendpoints.IPAllocationMethodDynamic && input.PrivateIPAddress != "" {
			return nil, fmt.Errorf("`private_ip_address` cannot be set when `private_ip_allocation_method` is `Dynamic`")
		}

		if input.PrivateIPAllocationMethod == inboundendpoints.IPAllocationMethodStatic && input.PrivateIPAddress == "" {
			return nil, fmt.Errorf("`private_ip_address` must be set when `private_ip_allocation_method` is `Static`")
		}

		if input.PrivateIPAllocationMethod != "" {
			output.PrivateIPAllocationMethod = &input.PrivateIPAllocationMethod
		}

		if input.PrivateIPAddress != "" {
			output.PrivateIPAddress = &input.PrivateIPAddress
		}

		output.Subnet = inboundendpoints.SubResource{
			Id: input.SubnetId,
		}

		outputList = append(outputList, output)
	}

	return &outputList, nil
}
