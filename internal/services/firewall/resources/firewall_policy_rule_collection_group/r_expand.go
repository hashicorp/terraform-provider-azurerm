// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package firewall_policy_rule_collection_group

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/firewallpolicyrulecollectiongroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandFirewallPolicyRuleCollectionApplication(input []any) []firewallpolicyrulecollectiongroups.FirewallPolicyRuleCollection {
	return expandFirewallPolicyFilterRuleCollection(input, expandFirewallPolicyRuleApplication)
}

func expandFirewallPolicyRuleCollectionNetwork(input []any) []firewallpolicyrulecollectiongroups.FirewallPolicyRuleCollection {
	return expandFirewallPolicyFilterRuleCollection(input, expandFirewallPolicyRuleNetwork)
}

func expandFirewallPolicyRuleCollectionNat(input []any) ([]firewallpolicyrulecollectiongroups.FirewallPolicyRuleCollection, error) {
	result := make([]firewallpolicyrulecollectiongroups.FirewallPolicyRuleCollection, 0)
	for _, e := range input {
		rule := e.(map[string]any)
		rules, err := expandFirewallPolicyRuleNat(rule["rule"].([]any))
		if err != nil {
			return nil, err
		}
		output := &firewallpolicyrulecollectiongroups.FirewallPolicyNatRuleCollection{
			Name:     pointer.To(rule["name"].(string)),
			Priority: pointer.To(int64(rule["priority"].(int))),
			Action: &firewallpolicyrulecollectiongroups.FirewallPolicyNatRuleCollectionAction{
				Type: pointer.ToEnum[firewallpolicyrulecollectiongroups.FirewallPolicyNatRuleCollectionActionType](rule["action"].(string)),
			},
			Rules: rules,
		}
		result = append(result, output)
	}
	return result, nil
}

func expandFirewallPolicyFilterRuleCollection(input []any, f func(input []any) *[]firewallpolicyrulecollectiongroups.FirewallPolicyRule) []firewallpolicyrulecollectiongroups.FirewallPolicyRuleCollection {
	result := make([]firewallpolicyrulecollectiongroups.FirewallPolicyRuleCollection, 0)
	for _, e := range input {
		rule := e.(map[string]any)
		output := &firewallpolicyrulecollectiongroups.FirewallPolicyFilterRuleCollection{
			Action: &firewallpolicyrulecollectiongroups.FirewallPolicyFilterRuleCollectionAction{
				Type: pointer.ToEnum[firewallpolicyrulecollectiongroups.FirewallPolicyFilterRuleCollectionActionType](rule["action"].(string)),
			},
			Name:     pointer.To(rule["name"].(string)),
			Priority: pointer.To(int64(rule["priority"].(int))),
			Rules:    f(rule["rule"].([]any)),
		}
		result = append(result, output)
	}
	return result
}

func expandFirewallPolicyRuleApplication(input []any) *[]firewallpolicyrulecollectiongroups.FirewallPolicyRule {
	result := make([]firewallpolicyrulecollectiongroups.FirewallPolicyRule, 0)
	for _, e := range input {
		condition := e.(map[string]any)
		var protocols []firewallpolicyrulecollectiongroups.FirewallPolicyRuleApplicationProtocol
		for _, p := range condition["protocols"].([]any) {
			proto := p.(map[string]any)
			protocols = append(protocols, firewallpolicyrulecollectiongroups.FirewallPolicyRuleApplicationProtocol{
				ProtocolType: pointer.ToEnum[firewallpolicyrulecollectiongroups.FirewallPolicyRuleApplicationProtocolType](proto["type"].(string)),
				Port:         pointer.To(int64(proto["port"].(int))),
			})
		}

		var httpHeader []firewallpolicyrulecollectiongroups.FirewallPolicyHTTPHeaderToInsert
		for _, h := range condition["http_headers"].([]any) {
			header := h.(map[string]any)
			httpHeader = append(httpHeader, firewallpolicyrulecollectiongroups.FirewallPolicyHTTPHeaderToInsert{
				HeaderName:  pointer.To(header["name"].(string)),
				HeaderValue: pointer.To(header["value"].(string)),
			})
		}

		output := &firewallpolicyrulecollectiongroups.ApplicationRule{
			Name:                 pointer.To(condition["name"].(string)),
			Description:          pointer.To(condition["description"].(string)),
			Protocols:            &protocols,
			HTTPHeadersToInsert:  &httpHeader,
			SourceAddresses:      pluginsdk.ExpandStringSlice(condition["source_addresses"].([]any)),
			SourceIPGroups:       pluginsdk.ExpandStringSlice(condition["source_ip_groups"].([]any)),
			DestinationAddresses: pluginsdk.ExpandStringSlice(condition["destination_addresses"].([]any)),
			TargetFqdns:          pluginsdk.ExpandStringSlice(condition["destination_fqdns"].([]any)),
			TargetURLs:           pluginsdk.ExpandStringSlice(condition["destination_urls"].([]any)),
			FqdnTags:             pluginsdk.ExpandStringSlice(condition["destination_fqdn_tags"].([]any)),
			TerminateTLS:         pointer.To(condition["terminate_tls"].(bool)),
			WebCategories:        pluginsdk.ExpandStringSlice(condition["web_categories"].([]any)),
		}
		result = append(result, output)
	}
	return &result
}

func expandFirewallPolicyRuleNetwork(input []any) *[]firewallpolicyrulecollectiongroups.FirewallPolicyRule {
	result := make([]firewallpolicyrulecollectiongroups.FirewallPolicyRule, 0)
	for _, e := range input {
		condition := e.(map[string]any)
		var protocols []firewallpolicyrulecollectiongroups.FirewallPolicyRuleNetworkProtocol
		for _, p := range condition["protocols"].([]any) {
			protocols = append(protocols, firewallpolicyrulecollectiongroups.FirewallPolicyRuleNetworkProtocol(p.(string)))
		}
		output := &firewallpolicyrulecollectiongroups.NetworkRule{
			Name:                 pointer.To(condition["name"].(string)),
			IPProtocols:          &protocols,
			SourceAddresses:      pluginsdk.ExpandStringSlice(condition["source_addresses"].([]any)),
			SourceIPGroups:       pluginsdk.ExpandStringSlice(condition["source_ip_groups"].([]any)),
			DestinationAddresses: pluginsdk.ExpandStringSlice(condition["destination_addresses"].([]any)),
			DestinationIPGroups:  pluginsdk.ExpandStringSlice(condition["destination_ip_groups"].([]any)),
			DestinationFqdns:     pluginsdk.ExpandStringSlice(condition["destination_fqdns"].([]any)),
			DestinationPorts:     pluginsdk.ExpandStringSlice(condition["destination_ports"].([]any)),
			Description:          pointer.To(condition["description"].(string)),
		}
		result = append(result, output)
	}
	return &result
}

func expandFirewallPolicyRuleNat(input []any) (*[]firewallpolicyrulecollectiongroups.FirewallPolicyRule, error) {
	result := make([]firewallpolicyrulecollectiongroups.FirewallPolicyRule, 0)
	for _, e := range input {
		condition := e.(map[string]any)
		var protocols []firewallpolicyrulecollectiongroups.FirewallPolicyRuleNetworkProtocol
		for _, p := range condition["protocols"].([]any) {
			protocols = append(protocols, firewallpolicyrulecollectiongroups.FirewallPolicyRuleNetworkProtocol(p.(string)))
		}

		// Exactly one of `translated_address` and `translated_fqdn` should be set.
		if condition["translated_address"].(string) != "" && condition["translated_fqdn"].(string) != "" {
			return nil, fmt.Errorf("can't specify both `translated_address` and `translated_fqdn` in rule %s", condition["name"].(string))
		}
		if condition["translated_address"].(string) == "" && condition["translated_fqdn"].(string) == "" {
			return nil, fmt.Errorf("should specify either `translated_address` or `translated_fqdn` in rule %s", condition["name"].(string))
		}
		output := &firewallpolicyrulecollectiongroups.NatRule{
			Name:                 pointer.To(condition["name"].(string)),
			IPProtocols:          &protocols,
			SourceAddresses:      pluginsdk.ExpandStringSlice(condition["source_addresses"].([]any)),
			SourceIPGroups:       pluginsdk.ExpandStringSlice(condition["source_ip_groups"].([]any)),
			DestinationAddresses: pointer.To([]string{condition["destination_address"].(string)}),
			DestinationPorts:     pluginsdk.ExpandStringSlice(condition["destination_ports"].([]any)),
			TranslatedPort:       pointer.To(strconv.Itoa(condition["translated_port"].(int))),
			Description:          pointer.To(condition["description"].(string)),
		}
		if condition["translated_address"].(string) != "" {
			output.TranslatedAddress = pointer.To(condition["translated_address"].(string))
		}
		if condition["translated_fqdn"].(string) != "" {
			output.TranslatedFqdn = pointer.To(condition["translated_fqdn"].(string))
		}
		result = append(result, output)
	}
	return &result, nil
}
