// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsforwardingrulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/virtualnetworklinks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PrivateDNSResolverVirtualNetworkLinDataSourcekModel struct {
	Name                   string            `tfschema:"name"`
	DnsForwardingRulesetId string            `tfschema:"dns_forwarding_ruleset_id"`
	Metadata               map[string]string `tfschema:"metadata"`
	VirtualNetworkId       string            `tfschema:"virtual_network_id"`
}

type PrivateDNSResolverVirtualNetworkLinkDataSource struct{}

var _ sdk.DataSource = PrivateDNSResolverVirtualNetworkLinkDataSource{}

func (r PrivateDNSResolverVirtualNetworkLinkDataSource) ResourceType() string {
	return "azurerm_private_dns_resolver_virtual_network_link"
}

func (r PrivateDNSResolverVirtualNetworkLinkDataSource) ModelObject() interface{} {
	return &PrivateDNSResolverVirtualNetworkLinDataSourcekModel{}
}

func (r PrivateDNSResolverVirtualNetworkLinkDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return virtualnetworklinks.ValidateVirtualNetworkLinkID
}

func (r PrivateDNSResolverVirtualNetworkLinkDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"dns_forwarding_ruleset_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: dnsforwardingrulesets.ValidateDnsForwardingRulesetID,
		},
	}
}

func (r PrivateDNSResolverVirtualNetworkLinkDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"virtual_network_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"metadata": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r PrivateDNSResolverVirtualNetworkLinkDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.VirtualNetworkLinksClient

			var state PrivateDNSResolverVirtualNetworkLinDataSourcekModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			dnsForwardingRulesetId, err := virtualnetworklinks.ParseDnsForwardingRulesetID(state.DnsForwardingRulesetId)
			if err != nil {
				return err
			}

			id := virtualnetworklinks.NewVirtualNetworkLinkID(
				dnsForwardingRulesetId.SubscriptionId,
				dnsForwardingRulesetId.ResourceGroupName,
				dnsForwardingRulesetId.DnsForwardingRulesetName,
				state.Name)
			resp, err := client.Get(ctx, id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state.DnsForwardingRulesetId = dnsforwardingrulesets.NewDnsForwardingRulesetID(id.SubscriptionId, id.ResourceGroupName, id.DnsForwardingRulesetName).ID()
			state.VirtualNetworkId = model.Properties.VirtualNetwork.Id

			if model.Properties.Metadata != nil {
				state.Metadata = *model.Properties.Metadata
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
