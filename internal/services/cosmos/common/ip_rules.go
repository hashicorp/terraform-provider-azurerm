// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2026-03-15/openapis"
)

// CosmosDBIpRulesToIpRangeFilterDataSource todo Remove for 4.0
func CosmosDBIpRulesToIpRangeFilterDataSource(ipRules *[]openapis.IPAddressOrRange) string {
	ipRangeFilter := make([]string, 0)
	if ipRules != nil {
		for _, ipRule := range *ipRules {
			ipRangeFilter = append(ipRangeFilter, *ipRule.IPAddressOrRange)
		}
	}

	return strings.Join(ipRangeFilter, ",")
}

func CosmosDBIpRulesToIpRangeFilter(ipRules *[]openapis.IPAddressOrRange) []string {
	ipRangeFilter := make([]string, 0)
	if ipRules != nil {
		for _, ipRule := range *ipRules {
			ipRangeFilter = append(ipRangeFilter, *ipRule.IPAddressOrRange)
		}
	}

	return ipRangeFilter
}

func CosmosDBIpRangeFilterToIpRules(ipRangeFilter []string) *[]openapis.IPAddressOrRange {
	ipRules := make([]openapis.IPAddressOrRange, 0)
	for _, ipRange := range ipRangeFilter {
		ipRules = append(ipRules, openapis.IPAddressOrRange{
			IPAddressOrRange: pointer.To(ipRange),
		})
	}

	return &ipRules
}

// CosmosDBIpRangeFilterToIpRulesThreePointOh todo Remove for 4.0
func CosmosDBIpRangeFilterToIpRulesThreePointOh(ipRangeFilter string) *[]openapis.IPAddressOrRange {
	ipRules := make([]openapis.IPAddressOrRange, 0)
	if len(ipRangeFilter) > 0 {
		for _, ipRange := range strings.Split(ipRangeFilter, ",") {
			ipRules = append(ipRules, openapis.IPAddressOrRange{
				IPAddressOrRange: pointer.To(ipRange),
			})
		}
	}

	return &ipRules
}
