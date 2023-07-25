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
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/inboundendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PrivateDNSResolverInboundEndpointDataSourceModel struct {
	Name                 string                 `tfschema:"name"`
	PrivateDNSResolverId string                 `tfschema:"private_dns_resolver_id"`
	IPConfigurations     []IPConfigurationModel `tfschema:"ip_configurations"`
	Location             string                 `tfschema:"location"`
	Tags                 map[string]string      `tfschema:"tags"`
}

type PrivateDNSResolverInboundEndpointDataSource struct{}

var _ sdk.DataSource = PrivateDNSResolverInboundEndpointDataSource{}

func (r PrivateDNSResolverInboundEndpointDataSource) ResourceType() string {
	return "azurerm_private_dns_resolver_inbound_endpoint"
}

func (r PrivateDNSResolverInboundEndpointDataSource) ModelObject() interface{} {
	return &PrivateDNSResolverInboundEndpointDataSourceModel{}
}

func (r PrivateDNSResolverInboundEndpointDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return inboundendpoints.ValidateInboundEndpointID
}

func (r PrivateDNSResolverInboundEndpointDataSource) Arguments() map[string]*pluginsdk.Schema {
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

func (r PrivateDNSResolverInboundEndpointDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"ip_configurations": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"private_ip_address": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"private_ip_allocation_method": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"location": commonschema.LocationComputed(),

		"tags": tags.SchemaDataSource(),
	}
}

func (r PrivateDNSResolverInboundEndpointDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.InboundEndpointsClient

			var state PrivateDNSResolverInboundEndpointDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			dnsForwardingRulesetId, err := inboundendpoints.ParseDnsResolverID(state.PrivateDNSResolverId)
			if err != nil {
				return err
			}

			id := inboundendpoints.NewInboundEndpointID(
				dnsForwardingRulesetId.SubscriptionId,
				dnsForwardingRulesetId.ResourceGroupName,
				dnsForwardingRulesetId.DnsResolverName,
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

			state.IPConfigurations = flattenIPConfigurationModel(&model.Properties.IPConfigurations)
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
