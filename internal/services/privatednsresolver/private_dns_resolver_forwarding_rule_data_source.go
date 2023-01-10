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
	Name                   string                           `tfschema:"name"`
	DnsForwardingRulesetId string                           `tfschema:"dns_forwarding_ruleset_id"`
	DomainName             string                           `tfschema:"domain_name"`
	Enabled                bool                             `tfschema:"enabled"`
	Metadata               map[string]string                `tfschema:"metadata"`
	TargetDnsServers       []TargetDnsServerDataSourceModel `tfschema:"target_dns_servers"`
}

type TargetDnsServerDataSourceModel struct {
	IPAddress string `tfschema:"ip_address"`
	Port      int64  `tfschema:"port"`
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

			var top int64 = 1
			resp, err := client.ListCompleteMatchingPredicate(ctx, *dnsForwardingRulesetId,
				forwardingrules.ListOperationOptions{Top: &top},
				forwardingrules.ForwardingRuleOperationPredicate{Name: &state.Name})
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", state.Name, err)
			}
			if len(resp.Items) != int(top) {
				return fmt.Errorf("retrieving %s: resource not found", state.Name)
			}

			model := resp.Items[0]
			id, err := forwardingrules.ParseForwardingRuleID(*model.Id)
			if err != nil {
				return err
			}

			state.DnsForwardingRulesetId = forwardingrules.NewDnsForwardingRulesetID(id.SubscriptionId, id.ResourceGroupName, id.DnsForwardingRulesetName).ID()
			properties := &model.Properties
			state.DomainName = properties.DomainName

			state.Enabled = false
			if properties.ForwardingRuleState != nil && *properties.ForwardingRuleState == forwardingrules.ForwardingRuleStateEnabled {
				state.Enabled = true
			}

			if properties.Metadata != nil {
				state.Metadata = *properties.Metadata
			}

			state.TargetDnsServers = flattenTargetDnsServerDataSouceModel(&properties.TargetDnsServers)

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}

func flattenTargetDnsServerDataSouceModel(inputList *[]forwardingrules.TargetDnsServer) []TargetDnsServerDataSourceModel {
	var outputList []TargetDnsServerDataSourceModel
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := TargetDnsServerDataSourceModel{
			IPAddress: input.IPAddress,
		}

		if input.Port != nil {
			output.Port = *input.Port
		}

		outputList = append(outputList, output)
	}

	return outputList
}
