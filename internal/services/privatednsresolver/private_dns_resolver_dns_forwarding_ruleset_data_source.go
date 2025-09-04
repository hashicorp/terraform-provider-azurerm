// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsforwardingrulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PrivateDNSResolverDnsForwardingRulesetDataSourceModel struct {
	Name                         string            `tfschema:"name"`
	ResourceGroupName            string            `tfschema:"resource_group_name"`
	DnsResolverOutboundEndpoints []string          `tfschema:"private_dns_resolver_outbound_endpoint_ids"`
	Location                     string            `tfschema:"location"`
	Tags                         map[string]string `tfschema:"tags"`
}

type PrivateDNSResolverDnsForwardingRulesetDataSource struct{}

var _ sdk.DataSource = PrivateDNSResolverDnsForwardingRulesetDataSource{}

func (r PrivateDNSResolverDnsForwardingRulesetDataSource) ResourceType() string {
	return "azurerm_private_dns_resolver_dns_forwarding_ruleset"
}

func (r PrivateDNSResolverDnsForwardingRulesetDataSource) ModelObject() interface{} {
	return &PrivateDNSResolverDnsForwardingRulesetDataSourceModel{}
}

func (r PrivateDNSResolverDnsForwardingRulesetDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dnsforwardingrulesets.ValidateDnsForwardingRulesetID
}

func (r PrivateDNSResolverDnsForwardingRulesetDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),
	}
}

func (r PrivateDNSResolverDnsForwardingRulesetDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"private_dns_resolver_outbound_endpoint_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"tags": tags.SchemaDataSource(),
	}
}

func (r PrivateDNSResolverDnsForwardingRulesetDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.DnsForwardingRulesetsClient

			var state PrivateDNSResolverDnsForwardingRulesetDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := dnsforwardingrulesets.NewDnsForwardingRulesetID(
				metadata.Client.Account.SubscriptionId, state.ResourceGroupName, state.Name)
			resp, err := client.Get(ctx, id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state.Location = location.Normalize(model.Location)
			state.DnsResolverOutboundEndpoints = flattenDnsResolverOutboundEndpoints(&model.Properties.DnsResolverOutboundEndpoints)

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
