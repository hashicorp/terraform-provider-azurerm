// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/networkanchors"
)

func FlattenDnsForwardingRules(input *[]networkanchors.DnsForwardingRule) []DnsForwardingRuleModel {
	output := make([]DnsForwardingRuleModel, 0)
	if input != nil {
		for _, item := range *input {
			output = append(output, DnsForwardingRuleModel{
				DomainNames:         item.DomainNames,
				ForwardingIPAddress: item.ForwardingIPAddress,
			})
		}
	}
	return output
}
