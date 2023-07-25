// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsresolvers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/outboundendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PrivateDNSResolverOutboundEndpointDataSourceModel struct {
	Name                 string            `tfschema:"name"`
	PrivateDNSResolverId string            `tfschema:"private_dns_resolver_id"`
	Location             string            `tfschema:"location"`
	SubnetId             string            `tfschema:"subnet_id"`
	Tags                 map[string]string `tfschema:"tags"`
}

type PrivateDNSResolverOutboundEndpointDataSource struct{}

var _ sdk.DataSource = PrivateDNSResolverOutboundEndpointDataSource{}

func (r PrivateDNSResolverOutboundEndpointDataSource) ResourceType() string {
	return "azurerm_private_dns_resolver_outbound_endpoint"
}

func (r PrivateDNSResolverOutboundEndpointDataSource) ModelObject() interface{} {
	return &PrivateDNSResolverOutboundEndpointDataSourceModel{}
}

func (r PrivateDNSResolverOutboundEndpointDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return outboundendpoints.ValidateOutboundEndpointID
}

func (r PrivateDNSResolverOutboundEndpointDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"private_dns_resolver_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: dnsresolvers.ValidateDnsResolverID,
		},
	}
}

func (r PrivateDNSResolverOutboundEndpointDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"subnet_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"location": commonschema.LocationComputed(),

		"tags": tags.SchemaDataSource(),
	}
}

func (r PrivateDNSResolverOutboundEndpointDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.OutboundEndpointsClient

			var state PrivateDNSResolverOutboundEndpointDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			privateDnsResolverId, err := outboundendpoints.ParseDnsResolverID(state.PrivateDNSResolverId)
			if err != nil {
				return err
			}

			id := outboundendpoints.NewOutboundEndpointID(
				privateDnsResolverId.SubscriptionId,
				privateDnsResolverId.ResourceGroupName,
				privateDnsResolverId.DnsResolverName,
				state.Name)
			resp, err := client.Get(ctx, id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state.Location = location.Normalize(model.Location)
			state.PrivateDNSResolverId = dnsresolvers.NewDnsResolverID(id.SubscriptionId, id.ResourceGroupName, id.DnsResolverName).ID()
			state.SubnetId = model.Properties.Subnet.Id

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
