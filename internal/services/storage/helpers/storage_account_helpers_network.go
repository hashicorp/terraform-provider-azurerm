// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func ExpandAccountNetworkRules(input []any, tenantId string) *storageaccounts.NetworkRuleSet {
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

	item := input[0].(map[string]any)
	return &storageaccounts.NetworkRuleSet{
		Bypass:              ExpandAccountNetworkRuleBypass(item["bypass"].(*pluginsdk.Set).List()),
		DefaultAction:       storageaccounts.DefaultAction(item["default_action"].(string)),
		IPRules:             ExpandAccountNetworkRuleIPRules(item["ip_rules"].(*pluginsdk.Set).List()),
		ResourceAccessRules: ExpandAccountNetworkRulePrivateLinkAccess(item["private_link_access"].([]any), tenantId),
		VirtualNetworkRules: ExpandAccountNetworkRuleVirtualNetworkRules(item["virtual_network_subnet_ids"].(*pluginsdk.Set).List()),
	}
}

func FlattenAccountNetworkRules(input *storageaccounts.NetworkRuleSet) []any {
	output := make([]any, 0)

	if input != nil {
		ipRules := FlattenAccountNetworkRuleIPRules(input.IPRules)
		privateLinkAccess := FlattenAccountNetworkRulePrivateLinkAccess(input.ResourceAccessRules)
		virtualNetworkRules := FlattenAccountNetworkRuleVirtualNetworkRules(input.VirtualNetworkRules)

		// ignore the default values
		usesDefaultValues := input.DefaultAction == storageaccounts.DefaultActionAllow && pointer.From(input.Bypass) == storageaccounts.BypassAzureServices
		hasOtherValues := len(ipRules) > 0 || len(privateLinkAccess) > 0 || len(virtualNetworkRules) > 0
		if usesDefaultValues && !hasOtherValues {
			return output
		}

		output = append(output, map[string]any{
			"bypass":                     pluginsdk.NewSet(pluginsdk.HashString, FlattenAccountNetworkRuleBypass(input.Bypass)),
			"default_action":             string(input.DefaultAction),
			"ip_rules":                   pluginsdk.NewSet(pluginsdk.HashString, ipRules),
			"private_link_access":        privateLinkAccess,
			"virtual_network_subnet_ids": pluginsdk.NewSet(pluginsdk.HashString, virtualNetworkRules),
		})
	}

	return output
}

func ExpandAccountNetworkRuleBypass(input []any) *storageaccounts.Bypass {
	if len(input) == 0 {
		return nil
	}

	output := make([]string, 0)
	for _, item := range input {
		output = append(output, item.(string))
	}
	return pointer.ToEnum[storageaccounts.Bypass](strings.Join(output, ", "))
}

func FlattenAccountNetworkRuleBypass(input *storageaccounts.Bypass) []any {
	output := make([]any, 0)

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

func ExpandAccountNetworkRuleIPRules(input []any) *[]storageaccounts.IPRule {
	output := make([]storageaccounts.IPRule, 0)
	for _, item := range input {
		output = append(output, storageaccounts.IPRule{
			Action: pointer.To(storageaccounts.ActionAllow),
			Value:  item.(string),
		})
	}
	return &output
}

func FlattenAccountNetworkRuleIPRules(input *[]storageaccounts.IPRule) []any {
	output := make([]any, 0)

	if input != nil {
		for _, item := range *input {
			output = append(output, item.Value)
		}
	}

	return output
}

func ExpandAccountNetworkRuleVirtualNetworkRules(input []any) *[]storageaccounts.VirtualNetworkRule {
	output := make([]storageaccounts.VirtualNetworkRule, 0)

	for _, item := range input {
		output = append(output, storageaccounts.VirtualNetworkRule{
			Id:     item.(string),
			Action: pointer.To(storageaccounts.ActionAllow),
		})
	}

	return &output
}

func FlattenAccountNetworkRuleVirtualNetworkRules(input *[]storageaccounts.VirtualNetworkRule) []any {
	output := make([]any, 0)

	if input != nil {
		for _, item := range *input {
			output = append(output, item.Id)
		}
	}

	return output
}

func ExpandAccountNetworkRulePrivateLinkAccess(input []any, tenantId string) *[]storageaccounts.ResourceAccessRule {
	output := make([]storageaccounts.ResourceAccessRule, 0)
	for _, raw := range input {
		item := raw.(map[string]any)
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

func FlattenAccountNetworkRulePrivateLinkAccess(input *[]storageaccounts.ResourceAccessRule) []any {
	output := make([]any, 0)

	if input != nil {
		for _, item := range *input {
			output = append(output, map[string]any{
				"endpoint_resource_id": pointer.From(item.ResourceId),
				"endpoint_tenant_id":   pointer.From(item.TenantId),
			})
		}
	}

	return output
}
