// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_resolver_forwarding_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsforwardingrulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/forwardingrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *forwardingrules.ForwardingRuleId, model *forwardingrules.ForwardingRule) error {
	state := PrivateDNSResolverForwardingRuleModel{
		Name:                   id.ForwardingRuleName,
		DnsForwardingRulesetId: dnsforwardingrulesets.NewDnsForwardingRulesetID(id.SubscriptionId, id.ResourceGroupName, id.DnsForwardingRulesetName).ID(),
	}

	if model != nil {
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
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func flattenTargetDnsServerModel(inputList *[]forwardingrules.TargetDnsServer) []TargetDnsServerModel {
	if inputList == nil {
		return []TargetDnsServerModel{}
	}

	outputList := make([]TargetDnsServerModel, 0, len(*inputList))
	for _, input := range *inputList {
		output := TargetDnsServerModel{
			IPAddress: input.IPAddress,
		}

		if input.Port != nil {
			output.Port = *input.Port
		}

		outputList = append(outputList, output)
	}

	return outputList
}
