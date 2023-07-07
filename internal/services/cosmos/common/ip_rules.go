// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// CosmosDBIpRulesToIpRangeFilterThreePointOh todo Remove for 4.0
func CosmosDBIpRulesToIpRangeFilterThreePointOh(ipRules *[]documentdb.IPAddressOrRange) string {
	ipRangeFilter := make([]string, 0)
	if ipRules != nil {
		for _, ipRule := range *ipRules {
			ipRangeFilter = append(ipRangeFilter, *ipRule.IPAddressOrRange)
		}
	}

	return strings.Join(ipRangeFilter, ",")
}

func CosmosDBIpRulesToIpRangeFilter(ipRules *[]documentdb.IPAddressOrRange) []string {
	ipRangeFilter := make([]string, 0)
	if ipRules != nil {
		for _, ipRule := range *ipRules {
			ipRangeFilter = append(ipRangeFilter, *ipRule.IPAddressOrRange)
		}
	}

	return ipRangeFilter
}

func CosmosDBIpRangeFilterToIpRules(ipRangeFilter []string) *[]documentdb.IPAddressOrRange {
	ipRules := make([]documentdb.IPAddressOrRange, 0)
	for _, ipRange := range ipRangeFilter {
		ipRules = append(ipRules, documentdb.IPAddressOrRange{
			IPAddressOrRange: utils.String(ipRange),
		})
	}

	return &ipRules
}

// CosmosDBIpRangeFilterToIpRulesThreePointOh todo Remove for 4.0
func CosmosDBIpRangeFilterToIpRulesThreePointOh(ipRangeFilter string) *[]documentdb.IPAddressOrRange {
	ipRules := make([]documentdb.IPAddressOrRange, 0)
	if len(ipRangeFilter) > 0 {
		for _, ipRange := range strings.Split(ipRangeFilter, ",") {
			ipRules = append(ipRules, documentdb.IPAddressOrRange{
				IPAddressOrRange: utils.String(ipRange),
			})
		}
	}

	return &ipRules
}
