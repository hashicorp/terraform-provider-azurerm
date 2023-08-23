// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsforwardingrulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/forwardingrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PrivateDNSResolverForwardingRuleDataSourceModel struct {
	Name                   string                 `tfschema:"name"`
	DnsForwardingRulesetId string                 `tfschema:"dns_forwarding_ruleset_id"`
	DomainName             string                 `tfschema:"domain_name"`
	Enabled                bool                   `tfschema:"enabled"`
	Metadata               map[string]string      `tfschema:"metadata"`
	TargetDnsServers       []TargetDnsServerModel `tfschema:"target_dns_servers"`
}

type PrivateDNSResolverForwardingRuleDataSource struct{}

var _ sdk.DataSource = PrivateDNSResolverForwardingRuleDataSource{}

func (r PrivateDNSResolverForwardingRuleDataSource) ResourceType() string {
	return "azurerm_private_dns_resolver_forwarding_rule"
}

func (r PrivateDNSResolverForwardingRuleDataSource) ModelObject() interface{} {
	return &PrivateDNSResolverForwardingRuleModel{}
}

func (r PrivateDNSResolverForwardingRuleDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return forwardingrules.ValidateForwardingRuleID
}

func (r PrivateDNSResolverForwardingRuleDataSource) Arguments() map[string]*pluginsdk.Schema {
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

func (r PrivateDNSResolverForwardingRuleDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"domain_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"target_dns_servers": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_address": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"port": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
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

func (r PrivateDNSResolverForwardingRuleDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.ForwardingRulesClient

			var state PrivateDNSResolverForwardingRuleDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			dnsForwardingRulesetId, err := forwardingrules.ParseDnsForwardingRulesetID(state.DnsForwardingRulesetId)
			if err != nil {
				return err
			}

			id := forwardingrules.NewForwardingRuleID(
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

			properties := &model.Properties
			state.DomainName = properties.DomainName

			state.Enabled = false
			if properties.ForwardingRuleState != nil && *properties.ForwardingRuleState == forwardingrules.ForwardingRuleStateEnabled {
				state.Enabled = true
			}

			if properties.Metadata != nil {
				state.Metadata = *properties.Metadata
			}

			state.TargetDnsServers = flattenTargetDnsServerModel(&properties.TargetDnsServers)

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
