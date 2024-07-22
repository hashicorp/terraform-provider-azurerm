// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandAccountNetworkRules(input []interface{}, tenantId string) *storageaccounts.NetworkRuleSet {
	if len(input) == 0 {
		// Default access is enabled when no network rules are set.
		return &storageaccounts.NetworkRuleSet{
			Bypass:              pointer.To(storageaccounts.BypassAzureServices),
			DefaultAction:       storageaccounts.DefaultActionAllow,
			IPRules:             &[]storageaccounts.IPRule{},
			ResourceAccessRules: &[]storageaccounts.ResourceAccessRule{},
			VirtualNetworkRules: &[]storageaccounts.VirtualNetworkRule{},
		}
	}

	item := input[0].(map[string]interface{})
	return &storageaccounts.NetworkRuleSet{
		Bypass:              expandAccountNetworkRuleBypass(item["bypass"].(*pluginsdk.Set).List()),
		DefaultAction:       storageaccounts.DefaultAction(item["default_action"].(string)),
		IPRules:             expandAccountNetworkRuleIPRules(item["ip_rules"].(*pluginsdk.Set).List()),
		ResourceAccessRules: expandAccountNetworkRulePrivateLinkAccess(item["private_link_access"].([]interface{}), tenantId),
		VirtualNetworkRules: expandAccountNetworkRuleVirtualNetworkRules(item["virtual_network_subnet_ids"].(*pluginsdk.Set).List()),
	}
}

func flattenAccountNetworkRules(input *storageaccounts.NetworkRuleSet) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		ipRules := flattenAccountNetworkRuleIPRules(input.IPRules)
		privateLinkAccess := flattenAccountNetworkRulePrivateLinkAccess(input.ResourceAccessRules)
		virtualNetworkRules := flattenAccountNetworkRuleVirtualNetworkRules(input.VirtualNetworkRules)

		// ignore the default values
		usesDefaultValues := input.DefaultAction == storageaccounts.DefaultActionAllow && pointer.From(input.Bypass) == storageaccounts.BypassAzureServices
		hasOtherValues := len(ipRules) > 0 || len(privateLinkAccess) > 0 || len(virtualNetworkRules) > 0
		if usesDefaultValues && !hasOtherValues {
			return output
		}

		output = append(output, map[string]interface{}{
			"bypass":                     pluginsdk.NewSet(pluginsdk.HashString, flattenAccountNetworkRuleBypass(input.Bypass)),
			"default_action":             string(input.DefaultAction),
			"ip_rules":                   pluginsdk.NewSet(pluginsdk.HashString, ipRules),
			"private_link_access":        privateLinkAccess,
			"virtual_network_subnet_ids": pluginsdk.NewSet(pluginsdk.HashString, virtualNetworkRules),
		})
	}

	return output
}

func expandAccountNetworkRuleBypass(input []interface{}) *storageaccounts.Bypass {
	if len(input) == 0 {
		return nil
	}

	output := make([]string, 0)
	for _, item := range input {
		output = append(output, item.(string))
	}
	return pointer.To(storageaccounts.Bypass(strings.Join(output, ", ")))
}

func flattenAccountNetworkRuleBypass(input *storageaccounts.Bypass) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		// Whilst this is an Enum it's actually a CSV containing the enum but its exposed as a regular string
		// as such we need to manually re-case this
		components := strings.Split(string(*input), ", ")

		vals := map[string]string{}
		for _, v := range storageaccounts.PossibleValuesForBypass() {
			vals[strings.ToLower(v)] = v
		}
		for _, item := range components {
			val, ok := vals[strings.ToLower(item)]
			if !ok {
				// otherwise append the direct casing if it's an unknown value
				val = item
			}

			output = append(output, val)
		}
	}

	return output
}

func expandAccountNetworkRuleIPRules(input []interface{}) *[]storageaccounts.IPRule {
	output := make([]storageaccounts.IPRule, 0)
	for _, item := range input {
		output = append(output, storageaccounts.IPRule{
			Action: pointer.To(storageaccounts.ActionAllow),
			Value:  item.(string),
		})
	}
	return &output
}

func flattenAccountNetworkRuleIPRules(input *[]storageaccounts.IPRule) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		for _, item := range *input {
			output = append(output, item.Value)
		}
	}

	return output
}

func expandAccountNetworkRuleVirtualNetworkRules(input []interface{}) *[]storageaccounts.VirtualNetworkRule {
	output := make([]storageaccounts.VirtualNetworkRule, 0)

	for _, item := range input {
		output = append(output, storageaccounts.VirtualNetworkRule{
			Id:     item.(string),
			Action: pointer.To(storageaccounts.ActionAllow),
		})
	}

	return &output
}

func flattenAccountNetworkRuleVirtualNetworkRules(input *[]storageaccounts.VirtualNetworkRule) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		for _, item := range *input {
			output = append(output, item.Id)
		}
	}

	return output
}

func expandAccountNetworkRulePrivateLinkAccess(input []interface{}, tenantId string) *[]storageaccounts.ResourceAccessRule {
	output := make([]storageaccounts.ResourceAccessRule, 0)
	for _, raw := range input {
		item := raw.(map[string]interface{})
		rule := storageaccounts.ResourceAccessRule{
			ResourceId: pointer.To(item["endpoint_resource_id"].(string)),
			TenantId:   pointer.To(tenantId),
		}
		if v := item["endpoint_tenant_id"].(string); v != "" {
			rule.TenantId = pointer.To(v)
		}
		output = append(output, rule)
	}

	return &output
}

func flattenAccountNetworkRulePrivateLinkAccess(input *[]storageaccounts.ResourceAccessRule) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		for _, item := range *input {
			output = append(output, map[string]interface{}{
				"endpoint_resource_id": pointer.From(item.ResourceId),
				"endpoint_tenant_id":   pointer.From(item.TenantId),
			})
		}
	}

	return output
}
