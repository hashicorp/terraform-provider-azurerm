// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
)

func flattenEncryption(input *registries.EncryptionProperty) []any {
	if input == nil || input.KeyVaultProperties == nil || input.Status == nil || *input.Status == registries.EncryptionStatusDisabled {
		return []any{}
	}

	return []any{
		map[string]any{
			"key_vault_key_id":   pointer.From(input.KeyVaultProperties.KeyIdentifier),
			"identity_client_id": pointer.From(input.KeyVaultProperties.Identity),
		},
	}
}

func flattenNetworkRuleSet(networkRuleSet *registries.NetworkRuleSet) []any {
	if networkRuleSet == nil {
		return []any{}
	}

	values := make(map[string]any)

	values["default_action"] = string(networkRuleSet.DefaultAction)

	ipRules := make([]any, 0)
	for _, ipRule := range *networkRuleSet.IPRules {
		value := make(map[string]any)
		value["action"] = string(*ipRule.Action)

		// When a /32 CIDR is passed as an ip rule, Azure will drop the /32 leading to the resource wanting to be re-created next run
		if !strings.Contains(ipRule.Value, "/") {
			ipRule.Value += "/32"
		}

		value["ip_range"] = ipRule.Value
		ipRules = append(ipRules, value)
	}

	values["ip_rule"] = ipRules

	return []any{values}
}

func flattenQuarantinePolicy(p *registries.Policies) bool {
	if p.QuarantinePolicy == nil {
		return false
	}

	return *p.QuarantinePolicy.Status == registries.PolicyStatusEnabled
}

func flattenExportPolicy(p *registries.Policies) bool {
	if p.ExportPolicy == nil {
		return false
	}

	return *p.ExportPolicy.Status == registries.ExportPolicyStatusEnabled
}

func flattenAadAuthAsArmPolicy(p *registries.Policies) bool {
	if p.AzureADAuthenticationAsArmPolicy == nil {
		return false
	}

	return *p.AzureADAuthenticationAsArmPolicy.Status == registries.AzureADAuthenticationAsArmPolicyStatusEnabled
}
