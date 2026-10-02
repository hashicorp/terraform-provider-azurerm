// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
)

func CosmosDBIpRulesToIpRangeFilter(ipRules *[]cosmosdb.IPAddressOrRange) []string {
	ipRangeFilter := make([]string, 0)
	if ipRules != nil {
		for _, ipRule := range *ipRules {
			ipRangeFilter = append(ipRangeFilter, *ipRule.IPAddressOrRange)
		}
	}

	return ipRangeFilter
}

func CosmosDBIpRangeFilterToIpRules(ipRangeFilter []string) *[]cosmosdb.IPAddressOrRange {
	ipRules := make([]cosmosdb.IPAddressOrRange, 0)
	for _, ipRange := range ipRangeFilter {
		ipRules = append(ipRules, cosmosdb.IPAddressOrRange{
			IPAddressOrRange: pointer.To(ipRange),
		})
	}

	return &ipRules
}
