// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// CosmosDBIpRulesToIpRangeFilterDataSource todo Remove for 4.0
func CosmosDBIpRulesToIpRangeFilterDataSource(ipRules *[]cosmosdb.IPAddressOrRange) string {
	ipRangeFilter := make([]string, 0)
	if ipRules != nil {
		for _, ipRule := range *ipRules {
			ipRangeFilter = append(ipRangeFilter, *ipRule.IPAddressOrRange)
		}
	}

	return strings.Join(ipRangeFilter, ",")
}

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
			IPAddressOrRange: utils.String(ipRange),
		})
	}

	return &ipRules
}

// CosmosDBIpRangeFilterToIpRulesThreePointOh todo Remove for 4.0
func CosmosDBIpRangeFilterToIpRulesThreePointOh(ipRangeFilter string) *[]cosmosdb.IPAddressOrRange {
	ipRules := make([]cosmosdb.IPAddressOrRange, 0)
	if len(ipRangeFilter) > 0 {
		for _, ipRange := range strings.Split(ipRangeFilter, ",") {
			ipRules = append(ipRules, cosmosdb.IPAddressOrRange{
				IPAddressOrRange: utils.String(ipRange),
			})
		}
	}

	return &ipRules
}
