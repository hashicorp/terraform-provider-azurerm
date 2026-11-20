// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_network_anchor

import (
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/networkanchors"
)

func expandDnsForwardingRules(dnsForwardingRules []DnsForwardingRuleModel) *[]networkanchors.DnsForwardingRule {
	results := make([]networkanchors.DnsForwardingRule, 0)
	for _, item := range dnsForwardingRules {
		results = append(results, networkanchors.DnsForwardingRule{
			DomainNames:         strings.Join(item.DomainNames, ","),
			ForwardingIPAddress: item.ForwardingIPAddress,
		})
	}
	return &results
}

func flattenDnsForwardingRules(dnsForwardingRules *[]networkanchors.DnsForwardingRule) []DnsForwardingRuleModel {
	results := make([]DnsForwardingRuleModel, 0)
	if dnsForwardingRules != nil {
		for _, item := range *dnsForwardingRules {
			results = append(results, DnsForwardingRuleModel{
				DomainNames:         flattenNetworkAnchorCommaSeparatedValues(item.DomainNames),
				ForwardingIPAddress: item.ForwardingIPAddress,
			})
		}
	}

	return results
}

func flattenNetworkAnchorCommaSeparatedValues(input string) []string {
	result := make([]string, 0)
	for value := range strings.SplitSeq(input, ",") {
		if value = strings.TrimSpace(value); value != "" {
			result = append(result, value)
		}
	}
	return result
}
