// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package firewall_policy_rule_collection_group

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/firewallpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/firewallpolicyrulecollectiongroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceFirewallPolicyRuleCollectionGroupSetFlatten(d *pluginsdk.ResourceData, id *firewallpolicyrulecollectiongroups.RuleCollectionGroupId, model *firewallpolicyrulecollectiongroups.FirewallPolicyRuleCollectionGroup) error {
	d.Set("name", id.RuleCollectionGroupName)
	d.Set("firewall_policy_id", firewallpolicies.NewFirewallPolicyID(id.SubscriptionId, id.ResourceGroupName, id.FirewallPolicyName).ID())

	if model != nil {
		if props := model.Properties; props != nil {
			d.Set("priority", props.Priority)

			applicationRuleCollections, networkRuleCollections, natRuleCollections, err := flattenFirewallPolicyRuleCollection(props.RuleCollections)
			if err != nil {
				return fmt.Errorf("flattening Firewall Policy Rule Collections: %+v", err)
			}

			if err := d.Set("application_rule_collection", applicationRuleCollections); err != nil {
				return fmt.Errorf("setting `application_rule_collection`: %+v", err)
			}
			if err := d.Set("network_rule_collection", networkRuleCollections); err != nil {
				return fmt.Errorf("setting `network_rule_collection`: %+v", err)
			}
			if err := d.Set("nat_rule_collection", natRuleCollections); err != nil {
				return fmt.Errorf("setting `nat_rule_collection`: %+v", err)
			}
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenFirewallPolicyRuleCollection(input *[]firewallpolicyrulecollectiongroups.FirewallPolicyRuleCollection) ([]any, []any, []any, error) {
	var (
		applicationRuleCollection = []any{}
		networkRuleCollection     = []any{}
		natRuleCollection         = []any{}
	)
	if input == nil {
		return applicationRuleCollection, networkRuleCollection, natRuleCollection, nil
	}

	for _, e := range *input {
		var result map[string]any

		switch rule := e.(type) {
		case firewallpolicyrulecollectiongroups.FirewallPolicyFilterRuleCollection:
			var action string
			if rule.Action != nil {
				action = pointer.FromEnum(rule.Action.Type)
			}

			result = map[string]any{
				"name":     pointer.From(rule.Name),
				"priority": pointer.From(rule.Priority),
				"action":   action,
			}

			if rule.Rules == nil || len(*rule.Rules) == 0 {
				continue
			}

			// Determine the rule type based on the first rule's type
			switch (*rule.Rules)[0].(type) {
			case firewallpolicyrulecollectiongroups.ApplicationRule:
				appRules, err := flattenFirewallPolicyRuleApplication(rule.Rules)
				if err != nil {
					return nil, nil, nil, err
				}
				result["rule"] = appRules

				applicationRuleCollection = append(applicationRuleCollection, result)

			case firewallpolicyrulecollectiongroups.FirewallPolicyRule:
				networkRules, err := flattenFirewallPolicyRuleNetwork(rule.Rules)
				if err != nil {
					return nil, nil, nil, err
				}
				result["rule"] = networkRules

				networkRuleCollection = append(networkRuleCollection, result)

			default:
				return nil, nil, nil, fmt.Errorf("unknown rule condition type %+v", (*rule.Rules)[0])
			}
		case firewallpolicyrulecollectiongroups.FirewallPolicyNatRuleCollection:
			var action string
			if rule.Action != nil {
				// doing this because we hardcode Dnat for https://github.com/Azure/azure-rest-api-specs/issues/9986
				if strings.EqualFold(pointer.FromEnum(rule.Action.Type), "Dnat") {
					action = "Dnat"
				} else {
					action = pointer.FromEnum(rule.Action.Type)
				}
			}

			rules, err := flattenFirewallPolicyRuleNat(rule.Rules)
			if err != nil {
				return nil, nil, nil, err
			}
			result = map[string]any{
				"name":     pointer.From(rule.Name),
				"priority": pointer.From(rule.Priority),
				"action":   action,
				"rule":     rules,
			}

			natRuleCollection = append(natRuleCollection, result)

		default:
			return nil, nil, nil, fmt.Errorf("unknown rule type %+v", rule)
		}
	}
	return applicationRuleCollection, networkRuleCollection, natRuleCollection, nil
}

func flattenFirewallPolicyRuleApplication(input *[]firewallpolicyrulecollectiongroups.FirewallPolicyRule) ([]any, error) {
	if input == nil {
		return []any{}, nil
	}
	output := make([]any, 0)
	for _, e := range *input {
		rule, ok := e.(firewallpolicyrulecollectiongroups.ApplicationRule)
		if !ok {
			return nil, fmt.Errorf("unexpected non-application rule: %+v", e)
		}

		protocols := make([]any, 0)
		if rule.Protocols != nil {
			for _, protocol := range *rule.Protocols {
				var port int
				if protocol.Port != nil {
					port = int(*protocol.Port)
				}
				protocols = append(protocols, map[string]any{
					"type": pointer.FromEnum(protocol.ProtocolType),
					"port": port,
				})
			}
		}

		httpHeaders := make([]any, 0)
		for _, header := range pointer.From(rule.HTTPHeadersToInsert) {
			httpHeaders = append(httpHeaders, map[string]any{
				"name":  pointer.From(header.HeaderName),
				"value": pointer.From(header.HeaderValue),
			})
		}

		output = append(output, map[string]any{
			"name":                  pointer.From(rule.Name),
			"description":           pointer.From(rule.Description),
			"protocols":             protocols,
			"http_headers":          httpHeaders,
			"source_addresses":      pluginsdk.FlattenSlice(rule.SourceAddresses),
			"source_ip_groups":      pluginsdk.FlattenSlice(rule.SourceIPGroups),
			"destination_addresses": pluginsdk.FlattenSlice(rule.DestinationAddresses),
			"destination_urls":      pluginsdk.FlattenSlice(rule.TargetURLs),
			"destination_fqdns":     pluginsdk.FlattenSlice(rule.TargetFqdns),
			"destination_fqdn_tags": pluginsdk.FlattenSlice(rule.FqdnTags),
			"terminate_tls":         pointer.From(rule.TerminateTLS),
			"web_categories":        pluginsdk.FlattenSlice(rule.WebCategories),
		})
	}

	return output, nil
}

func flattenFirewallPolicyRuleNetwork(input *[]firewallpolicyrulecollectiongroups.FirewallPolicyRule) ([]any, error) {
	if input == nil {
		return []any{}, nil
	}
	output := make([]any, 0)
	for _, e := range *input {
		rule, ok := e.(firewallpolicyrulecollectiongroups.NetworkRule)
		if !ok {
			return nil, fmt.Errorf("unexpected non-network rule: %+v", e)
		}

		protocols := make([]any, 0)
		if rule.IPProtocols != nil {
			for _, protocol := range *rule.IPProtocols {
				protocols = append(protocols, string(protocol))
			}
		}

		output = append(output, map[string]any{
			"name":                  pointer.From(rule.Name),
			"protocols":             protocols,
			"source_addresses":      pluginsdk.FlattenSlice(rule.SourceAddresses),
			"source_ip_groups":      pluginsdk.FlattenSlice(rule.SourceIPGroups),
			"destination_addresses": pluginsdk.FlattenSlice(rule.DestinationAddresses),
			"destination_ip_groups": pluginsdk.FlattenSlice(rule.DestinationIPGroups),
			"destination_fqdns":     pluginsdk.FlattenSlice(rule.DestinationFqdns),
			"destination_ports":     pluginsdk.FlattenSlice(rule.DestinationPorts),
			"description":           pointer.From(rule.Description),
		})
	}
	return output, nil
}

func flattenFirewallPolicyRuleNat(input *[]firewallpolicyrulecollectiongroups.FirewallPolicyRule) ([]any, error) {
	if input == nil {
		return []any{}, nil
	}
	output := make([]any, 0)
	for _, e := range *input {
		rule, ok := e.(firewallpolicyrulecollectiongroups.NatRule)
		if !ok {
			return nil, fmt.Errorf("unexpected non-nat rule: %+v", e)
		}

		protocols := make([]any, 0)
		if rule.IPProtocols != nil {
			for _, protocol := range *rule.IPProtocols {
				protocols = append(protocols, string(protocol))
			}
		}
		destinationAddr := ""
		if rule.DestinationAddresses != nil && len(*rule.DestinationAddresses) != 0 {
			destinationAddr = (*rule.DestinationAddresses)[0]
		}

		translatedPort := 0
		if rule.TranslatedPort != nil {
			port, err := strconv.Atoi(*rule.TranslatedPort)
			if err != nil {
				return nil, fmt.Errorf(`the "translatedPort" property is not a valid integer (%s)`, *rule.TranslatedPort)
			}
			translatedPort = port
		}

		output = append(output, map[string]any{
			"name":                pointer.From(rule.Name),
			"protocols":           protocols,
			"source_addresses":    pluginsdk.FlattenSlice(rule.SourceAddresses),
			"source_ip_groups":    pluginsdk.FlattenSlice(rule.SourceIPGroups),
			"destination_address": destinationAddr,
			"destination_ports":   pluginsdk.FlattenSlice(rule.DestinationPorts),
			"translated_address":  pointer.From(rule.TranslatedAddress),
			"translated_port":     translatedPort,
			"translated_fqdn":     pointer.From(rule.TranslatedFqdn),
			"description":         pointer.From(rule.Description),
		})
	}
	return output, nil
}
