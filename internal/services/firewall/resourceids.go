// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package firewall

// Firewall Policy
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FirewallApplicationRuleCollection -id=/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/azureFirewalls/myfirewall/applicationRuleCollections/applicationRuleCollection1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FirewallNatRuleCollection -id=/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/azureFirewalls/myfirewall/natRuleCollections/natRuleCollection1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FirewallNetworkRuleCollection -id=/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/azureFirewalls/myfirewall/networkRuleCollections/networkRuleCollection1
